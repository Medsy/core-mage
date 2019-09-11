package targets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/LUSHDigital/core-mage/targets/chart"
	"github.com/LUSHDigital/core-mage/targets/compose"

	"github.com/joho/godotenv"
)

var (
	// InfraDir represents the directory used for infrastructure files.
	InfraDir = "infra"

	// GCPStageChartFile configures the file path to the gcp chart for stage.
	GCPStageChartFile = path.Join(InfraDir, "stage.gcp.yaml")

	// GCPProdChartFile configures the file path to the gcp chart for prod.
	GCPProdChartFile = path.Join(InfraDir, "prod.gcp.yaml")
)

type chartData struct {
	Name string
	Type string
}

func writeInfraDir() error {
	return os.MkdirAll(InfraDir, os.ModePerm)
}

func writeStageChart() error {
	return writeChart(GCPStageChartFile)
}

func writeProdChart() error {
	return writeChart(GCPProdChartFile)
}

func writeDotEnvFiles() error {
	if err := writeDotEnv(); err != nil {
		return err
	}
	if err := writeDotEnvLocal(); err != nil {
		return err
	}
	if err := writeDotEnvDev(); err != nil {
		return err
	}
	if err := writeDotEnvTest(); err != nil {
		return err
	}
	return nil
}

func writeDotEnv() error {
	var vars = make(map[string]string)
	for _, dep := range DockerComposeDevDependencies {
		k, v := compose.Services.EnvFor(LocalHost, dep)
		vars[k] = v
	}
	if DockerRunImage == DockerRunImageMigrations {
		vars[MigrationsURLEnvVar] = MigrationsURLLocal
	}
	return writeEnvFile(path.Join(InfraDir, ".env"), vars)
}

func writeDotEnvDev() error {
	var vars = make(map[string]string)
	for _, dep := range DockerComposeDevDependencies {
		k, v := compose.Services.EnvFor(dep, dep)
		vars[k] = v
	}
	if DockerRunImage == DockerRunImageMigrations {
		vars[MigrationsURLEnvVar] = MigrationsURLDev
	}
	return writeEnvFile(path.Join(InfraDir, "dev.env"), vars)
}

func writeDotEnvTest() error {
	var vars = make(map[string]string)
	for _, dep := range DockerComposeTestDependencies {
		k, v := compose.Services.EnvFor(dep, dep)
		vars[k] = v
	}
	if DockerRunImage == DockerRunImageMigrations {
		vars[MigrationsURLEnvVar] = MigrationsURLTest
	}
	return writeEnvFile(path.Join(InfraDir, "test.env"), vars)
}

func writeDotEnvLocal() error {
	var vars = make(map[string]string)
	return writeEnvFile(path.Join(InfraDir, "local.env"), vars)
}

func writeEnvFile(filename string, vars map[string]string) error {
	filevars, err := godotenv.Read(filename)
	if err != nil {
		fmt.Printf("env file %q does not exist: skipping...\n", filename)
	}
	var allvars = make(map[string]string)
	for k, v := range vars {
		allvars[strings.ToUpper(k)] = v
	}
	for k, v := range filevars {
		allvars[strings.ToUpper(k)] = v
	}
	buf := bytes.NewBuffer(nil)
	for k, v := range allvars {
		fmt.Fprintf(buf, "%s=%s\n", k, v)
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0664)
}

func writeChart(filename string) error {
	_, err := os.Stat(filename)
	if err != nil && os.IsExist(err) {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if os.IsNotExist(err) {
		chartVersion := "9.0.0-stable"
		pullPolicy := "Always"
		replicas := 3
		v9f := chart.V9File{
			ChartVersion: &chartVersion,
			ServiceName:  &ProjectName,
			ServiceType:  &ProjectType,
			Replicas:     &replicas,
			PullPolicy:   &pullPolicy,
		}
		if _, err := v9f.WriteTo(buf); err != nil {
			return err
		}
	} else {
		vf, err := chart.ReadVersionFile(filename)
		if err != nil {
			return err
		}
		switch vf.Version() {
		case chart.V9:
			v9f, err := chart.ReadV9File(filename)
			if err != nil {
				return err
			}
			_, err = v9f.WriteTo(buf)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown chart version: %q", vf.ChartVersion)
		}
	}
	return ioutil.WriteFile(filename, []byte(strings.TrimSpace(buf.String())), 0664)
}
