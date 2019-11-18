// Copyright 2019 The FreeIPA Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
    "bytes"
    "fmt"
    "github.com/GeertJohan/go.rice"
    "github.com/briantopping/freeipa-operator/api/v1alpha1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/kubernetes/scheme"
    "strings"
    "text/template"
)

type RenderContext struct {
    Key     string
    Cluster *v1alpha1.IpaCluster
    Server  *v1alpha1.Server
}

// Creates a list of objects from a YAML template. These objects are introspected and applied by the caller
func (ctx *RenderContext) ProcessTemplate() ([]runtime.Object, error) {
    templateBox, err := rice.FindBox(".")
    if err != nil {
        return nil, fmt.Errorf("could not open templates box, error is %e", err)
    }
    // get file contents as string
    templateString, err := templateBox.String(ctx.Key)
    if err != nil {
        return nil, fmt.Errorf("could not open template, error is %e", err)
    }
    t, err := template.New("template").Funcs(funcMaps).Parse(templateString)
    if err != nil {
        return nil, fmt.Errorf("could not create template, error is %e", err)
    }
    templateBuf := &bytes.Buffer{}
    err = t.Execute(templateBuf, ctx)
    if err != nil {
        return nil, fmt.Errorf("could not execute template, error is %e", err)
    }

    files := strings.Split(templateBuf.String(), "---")
    decode := scheme.Codecs.UniversalDeserializer().Decode
    var objs []runtime.Object
    for _, f := range files {
        if f == "\n" || f == "" {
            // ignore empty cases
            continue
        }

        obj, _, e := decode([]byte(f), nil, nil)

        if e != nil {
            return nil, fmt.Errorf("could not decode template, error is %e", e)
        }
        objs = append(objs, obj)
    }
    return objs, nil
}

var funcMaps = template.FuncMap{
    "seq": func(c int) []interface{} {
        return make([]interface{}, c)
    },
}

