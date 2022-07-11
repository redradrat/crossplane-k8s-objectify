package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha1"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	k8syaml "sigs.k8s.io/yaml"
)

var (
	input  string
	output string
)

type IntObject struct {
	v1.TypeMeta `json:",inline"`
	Spec        v1alpha1.ObjectSpec `json:"spec"`
}

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println("error constructing logger")
		os.Exit(1)
	}
	flag.StringVarP(&input, "input", "i", "", "The file to convert to crossplane k8s Objects")
	flag.StringVarP(&output, "output", "o", "objects.yaml", "The file to write the converted output to")
	flag.Parse()

	parsedFile, err := parseInput(input)
	if err != nil {
		log.Sugar().Fatalf("issue encountered while parsing input file: %v", err)
	}

	var out []byte
	for _, f := range parsedFile {

		obj := IntObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Object",
				APIVersion: "kubernetes.crossplane.io/v1alpha1",
			},
			Spec: v1alpha1.ObjectSpec{
				ForProvider: v1alpha1.ObjectParameters{
					Manifest: runtime.RawExtension{Raw: f},
				},
			},
		}

		objDataJson, err := json.Marshal(&obj)
		if err != nil {
			log.Sugar().Fatalf("unable to marshal parsed data into target file: %v", err)
		}
		objDataYaml, err := k8syaml.JSONToYAML(objDataJson)
		if err != nil {
			log.Sugar().Fatalf("unable to marshal rendered data to YAML: %v", err)
		}
		out = append(out, objDataYaml...)
		out = append(out, []byte("---\n")...)
	}
	if err := ioutil.WriteFile(output, out, 0664); err != nil {
		log.Sugar().Fatalf("unable to write target file: %v", err)
	}

	return
}

type T map[string]interface{}

func parseInput(i string) ([][]byte, error) {
	var out [][]byte
	yfile, err := ioutil.ReadFile(i)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(yfile))
	for {
		obj := &T{}
		err := decoder.Decode(&obj)

		if obj == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}

		marshalledJson, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		out = append(out, marshalledJson)
	}

	return out, nil
}
