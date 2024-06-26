/*
Copyright 2024 KubeAGI.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package arctl

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	basev1alpha1 "github.com/kubeagi/arcadia/api/base/v1alpha1"
	evalv1alpha1 "github.com/kubeagi/arcadia/api/evaluation/v1alpha1"
	pkgclient "github.com/kubeagi/arcadia/apiserver/pkg/client"
	downloadutil "github.com/kubeagi/arcadia/pkg/arctl/download"
	"github.com/kubeagi/arcadia/pkg/config"
	"github.com/kubeagi/arcadia/pkg/evaluation"
	pkgutils "github.com/kubeagi/arcadia/pkg/utils"
)

func NewEvalCmd(home *string, namespace *string) *cobra.Command {
	var appName string

	cmd := &cobra.Command{
		Use:   "eval",
		Short: "Manage evaluations",
	}

	cmd.PersistentFlags().StringVar(&appName, "application", "", "The application to be evaluated")
	if err := cmd.MarkPersistentFlagRequired("application"); err != nil {
		panic(err)
	}

	cmd.AddCommand(EvalGenTestDataset(home, namespace, &appName))
	cmd.AddCommand(NewRAGDownloadCmd(namespace))
	return cmd
}

func NewRAGDownloadCmd(namespace *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "download all test files in the RAG selected dataset",
		Long: `Downloading files locally from minio.
Example:

arctl -narcadia eval --rag=<rag-name>
`,
		Run: func(cmd *cobra.Command, args []string) {
			kubeClient, err := pkgclient.GetClient(nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to connect cluster error %s\n", err)
				os.Exit(1)
			}
			ragName, _ := cmd.Flags().GetString("rag")
			dir, _ := cmd.Flags().GetString("dir")
			systemConfNamespace, _ := cmd.Flags().GetString("system-conf-namespace")
			systemConfName, _ := cmd.Flags().GetString("system-conf-name")
			rag := &evalv1alpha1.RAG{}
			err = kubeClient.Get(context.TODO(), types.NamespacedName{Namespace: *namespace, Name: ragName}, rag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get rag %s error %s\n", ragName, err)
				os.Exit(1)
			}

			download(cmd.Context(), rag, kubeClient, dir, systemConfName, systemConfNamespace)
		},
	}
	cmd.Flags().String("rag", "", "rag name")
	cmd.Flags().String("dir", ".", "specify the file download directory")
	cmd.Flags().String("system-conf-namespace", "", "the namespace where the system configuration of the arcadia service is located.")
	cmd.Flags().String("system-conf-name", "arcadia-config", "the system configuration name of the arcadia service.")
	_ = cmd.MarkFlagRequired("rag")
	return cmd
}

func EvalGenTestDataset(home *string, namespace *string, appName *string) *cobra.Command {
	var (
		inputDir           string
		questionColumn     string
		groundTruthsColumn string
		outputMethod       string
		outputDir          string
		mergeFileName      string
		merge              bool
		maxContentLength   int
	)

	cmd := &cobra.Command{
		Use:   "gen_test_dataset",
		Short: "Generate a test dataset for evaluation",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if outputDir == "" {
				outputDir = *home
			}

			// init kubeclient
			kubeClient, err := pkgclient.GetClient(nil)
			if err != nil {
				return err
			}

			// read files
			app := &basev1alpha1.Application{}
			err = kubeClient.Get(ctx, types.NamespacedName{Namespace: *namespace, Name: *appName}, app)
			if err != nil {
				return err
			}

			var (
				csvWriter   *csv.Writer
				writeHeader = true
			)
			if merge {
				mergeFile, err := os.OpenFile(mergeFileName, os.O_CREATE|os.O_RDWR, 0744)
				if err != nil {
					return err
				}
				defer mergeFile.Close()
				csvWriter = csv.NewWriter(mergeFile)
			}
			// read files from input directory
			err = filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() || filepath.Ext(d.Name()) != ".csv" || strings.HasPrefix(d.Name(), "ragas-") {
					return nil
				}
				var output evaluation.Output
				switch outputMethod {
				case "csv":
					csvOutput := &evaluation.CSVOutput{
						W: csvWriter,
					}
					if !merge {
						outputCSVFile, err := os.Create(strings.Replace(path, d.Name(), fmt.Sprintf("ragas-%s", d.Name()), 1))
						if err != nil {
							return err
						}
						defer outputCSVFile.Close()
						csvOutput.W = csv.NewWriter(outputCSVFile)
					}
					defer csvOutput.W.Flush()
					output = csvOutput
				default:
					output = &evaluation.PrintOutput{}
				}
				// read file from dataset
				err = GenDatasetOnSingleFile(ctx, kubeClient, app,
					path,
					evaluation.WithQuestionColumn(questionColumn),
					evaluation.WithGroundTruthsColumn(groundTruthsColumn),
					evaluation.WithOutput(output),
					evaluation.WithWriteHeader(!merge || writeHeader),
					evaluation.WithMaxContextLength(maxContentLength),
				)
				if err != nil {
					return err
				}
				writeHeader = false
				return nil
			})
			return err
		},
	}

	cmd.Flags().StringVar(&inputDir, "input-dir", "", "The input directory where to load original dataset files")
	if err := cmd.MarkFlagRequired("input-dir"); err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&questionColumn, "question-column", "q", "The column name which provides questions")
	cmd.Flags().StringVar(&groundTruthsColumn, "ground-truths-column", "a", "The column name which provides the answers")
	cmd.Flags().StringVar(&outputMethod, "output", "", "The way to output the generated dataset rows.We support two ways: \n - stdout: print row \n - csv: save row to csv file")
	cmd.Flags().BoolVar(&merge, "merge", false, "Whether to merge all generated test data into a single file")
	cmd.Flags().StringVar(&mergeFileName, "merge-file", "ragas.csv", "name of the merged document")
	cmd.Flags().IntVar(&maxContentLength, "max-context-length", 512, "The maximum length of the context")

	return cmd
}

func GenDatasetOnSingleFile(ctx context.Context, kubeClient client.Client, app *basev1alpha1.Application, file string, genOpts ...evaluation.GenOptions) error {
	klog.V(3).Infof("GenDatasetOnSingleFile for application %s in namespace %s from file %s", app.Name, app.Namespace, file)
	// read file content
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// init evaluation dataset generator
	generator, err := evaluation.NewRagasDatasetGenerator(ctx, kubeClient, app, genOpts...)
	if err != nil {
		return err
	}

	// generate test dataset
	err = generator.Generate(
		ctx,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetCustomDefineDatasource(
	ctx context.Context,
	kubeClient client.Client,
	datasourceName, datasourceNamespace string,
) (*basev1alpha1.Datasource, error) {
	datasource := &basev1alpha1.Datasource{}
	err := kubeClient.Get(ctx, types.NamespacedName{Namespace: datasourceNamespace, Name: datasourceName}, datasource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get datasource %s/%s error %s\n", datasourceNamespace, datasourceName, err)
		return nil, err
	}
	return datasource, nil
}

var (
	once             sync.Once
	systemDatasource *basev1alpha1.Datasource
	systemError      error
)

func SysatemDatasource(ctx context.Context, kubeClient client.Client) (*basev1alpha1.Datasource, error) {
	once.Do(func() {
		systemDatasource, systemError = config.GetSystemDatasource(ctx)
	})
	return systemDatasource, systemError
}

func parseOptions(
	ctx context.Context,
	kubeClient client.Client,
	datasource *basev1alpha1.Datasource,
) []downloadutil.DownloadOptionFunc {
	options := make([]downloadutil.DownloadOptionFunc, 0)
	options = append(options, downloadutil.WithEndpoint(datasource.Spec.Endpoint.URL),
		downloadutil.WithSecure(datasource.Spec.Endpoint.Insecure))

	if as := datasource.Spec.Endpoint.AuthSecret; as != nil {
		secret := &v1.Secret{}
		ns := datasource.Namespace
		if as.Namespace != nil {
			ns = *as.Namespace
		}
		err := kubeClient.Get(ctx, types.NamespacedName{Namespace: ns, Name: as.Name}, secret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get auth secret %s error %s\n", as.Name, err)
			return nil
		}
		pwd := secret.Data["rootPassword"]
		user := secret.Data["rootUser"]
		options = append(options, downloadutil.WithAccessKey(string(user)), downloadutil.WithSecretKey(string(pwd)))
	}

	return options
}

func fromDatasource(ctx context.Context, kubeClient client.Client, dsName, dsNamespace string, clientCache map[string]*downloadutil.Download) error {
	key := fmt.Sprintf("%s/%s", dsNamespace, dsName)
	_, ok := clientCache[key]
	if ok {
		return nil
	}
	datasource, err := GetCustomDefineDatasource(ctx, kubeClient, dsName, dsNamespace)
	if err != nil {
		return err
	}
	if datasource.Spec.OSS == nil {
		return errors.New("the datasource is not configured with bucket information")
	}

	options := parseOptions(ctx, kubeClient, datasource)
	options = append(options, downloadutil.WithBucket(datasource.Spec.OSS.Bucket))
	downloader := downloadutil.NewDownloader(options...)
	clientCache[key] = downloader
	return nil
}

const (
	systemDatasourceKey = "system-datasource"
)

func fromVersionedDataset(ctx context.Context, kubeClient client.Client, bucket string, clientCache map[string]*downloadutil.Download) error {
	datasource, err := SysatemDatasource(ctx, kubeClient)
	if err != nil {
		return err
	}
	if datasource.Spec.OSS == nil {
		return errors.New("the system datasource is not configured with bucket information")
	}
	_, ok := clientCache[systemDatasourceKey]
	if ok {
		return nil
	}

	options := parseOptions(ctx, kubeClient, datasource)
	options = append(options, downloadutil.WithBucket(bucket))

	downloader := downloadutil.NewDownloader(options...)
	clientCache[systemDatasourceKey] = downloader
	return nil
}

func download(
	ctx context.Context,
	rag *evalv1alpha1.RAG,
	kubeClient client.Client,
	baseDir, systemConfName, systemConfNamespace string,
) {
	curNamespace := rag.Namespace
	clientCache := make(map[string]*downloadutil.Download)

	os.Setenv(config.EnvConfigKey, systemConfName)
	os.Setenv(pkgutils.EnvNamespaceKey, systemConfNamespace)
	defer func() {
		os.Unsetenv(config.EnvConfigKey)
		os.Unsetenv(pkgutils.EnvNamespaceKey)
	}()

	for _, source := range rag.Spec.Datasets {
		if source.Source.Kind != "Datasource" && source.Source.Kind != "VersionedDataset" {
			fmt.Fprintf(os.Stderr, "warning: only support Datasource, VersioneddataSet to get data, the current fill in the kind is %s", source.Source.Kind)
			continue
		}

		key := ""
		ns := curNamespace
		if source.Source.Kind == "Datasource" {
			if source.Source.Namespace != nil {
				ns = *source.Source.Namespace
			}
			key = fmt.Sprintf("%s/%s", ns, source.Source.Name)
			if err := fromDatasource(ctx, kubeClient, source.Source.Name, ns, clientCache); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get datasource %s error %s", source.Source.Name, err)
				continue
			}
		}
		if source.Source.Kind == "VersionedDataset" {
			if source.Source.Namespace != nil {
				ns = *source.Source.Namespace
			}
			key = systemDatasourceKey
			if err := fromVersionedDataset(ctx, kubeClient, ns, clientCache); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get system datasource error %s", err)
				continue
			}
		}

		downloader := clientCache[key]
		for _, f := range source.Files {
			if err := downloader.Download(ctx, downloadutil.WithSrcFile(f), downloadutil.WithDstFile(filepath.Join(baseDir, f))); err != nil {
				fmt.Fprintf(os.Stderr, "datasource %s failed to download file error %s\n", key, err)
			}
		}
		clientCache[key] = downloader
	}
}
