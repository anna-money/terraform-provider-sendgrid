package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/octoenergy/terraform-provider-sendgrid/sendgrid"
)

const (
	providerName   = "sendgrid"
	indexFile      = "index.md.tpl"
	docFile        = "doc.md.tpl"
	docRoot        = "../docs"
	dataSourcesStr = "Data Sources"
	forceNewStr    = ", ForceNew"
)

func main() {
	provider := sendgrid.Provider()
	vProvider := runtime.FuncForPC(reflect.ValueOf(sendgrid.Provider).Pointer())

	fname, _ := vProvider.FileLine(0)
	fpath := filepath.Dir(fname)
	log.Printf("generating doc from: %s\n", fpath)

	// document for DataSources
	for k, v := range provider.DataSourcesMap {
		genDoc("data_source", "data-sources", fpath, k, v)
	}

	// document for Resources
	for k, v := range provider.ResourcesMap {
		genDoc("resource", "resources", fpath, k, v)
	}

	// document for Index
	genIdx(fpath)
}

// genIdx generating index for resource.
func genIdx(fpath string) { //nolint:cyclop,funlen
	type Index struct {
		Name          string
		NameShort     string
		ResType       string
		ResTypeFolder string
		Resources     [][]string
	}

	var resources string

	var dataSources []Index

	var sources []Index

	fname := "provider.go"
	log.Printf("[START]get description from file: %s\n", fname)

	description, err := getFileDescription(fmt.Sprintf("%s/%s", fpath, fname))
	if err != nil {
		log.Printf("[SKIP!]get description failed, skip: %s", err)

		return
	}

	description = strings.TrimSpace(description)
	if description == "" {
		log.Printf("[SKIP!]description empty, skip: %s\n", fname)

		return
	}

	pos := strings.Index(description, "Resources List\n")
	if pos != -1 {
		resources = strings.TrimSpace(description[pos+16:])
	} else {
		log.Printf("[SKIP!]resource list missing, skip: %s\n", fname)

		return
	}

	index := Index{}

	for _, v := range strings.Split(resources, "\n") {
		vv := strings.TrimSpace(v)
		if vv == "" {
			continue
		}

		if strings.HasPrefix(v, "  ") { //nolint:nestif
			if index.Name == "" {
				log.Printf("[FAIL!]no resource name found: %s", v)

				return
			}

			index.Resources = append(index.Resources, []string{vv, vv[len(providerName)+1:]})
		} else {
			if index.Name != "" {
				if index.Name == dataSourcesStr {
					dataSources = append(dataSources, index)
				} else {
					sources = append(sources, index)
				}
			}
			vvv := ""
			resType := "datasource"
			resTypeFolder := "data-sources"
			if vv != dataSourcesStr {
				resType = "resource"
				resTypeFolder = "resources"
				vs := strings.Split(vv, " ")
				vvv = strings.ToLower(strings.Join(vs[:len(vs)-1], "-"))
			}
			index = Index{
				Name:          vv,
				NameShort:     vvv,
				ResType:       resType,
				ResTypeFolder: resTypeFolder,
				Resources:     [][]string{},
			}
		}
	}

	if index.Name != "" {
		if index.Name == dataSourcesStr {
			dataSources = append(dataSources, index)
		} else {
			sources = append(sources, index)
		}
	}

	dataSources = append(dataSources, sources...)
	data := map[string]interface{}{
		"datasource": dataSources,
	}

	fname = fmt.Sprintf("%s/index.md", docRoot)

	fd, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644) //nolint:gomnd
	if err != nil {
		log.Printf("[FAIL!]open file %s failed: %s", fname, err)

		return
	}

	defer func() {
		if e := fd.Close(); e != nil {
			log.Printf("[FAIL!]close file %s failed: %s", fname, e)
		}
	}()

	idxTPL, err := os.ReadFile(indexFile)
	if err != nil {
		log.Printf("[FAIL!]open file %s failed: %s", indexFile, err)

		return
	}

	t := template.Must(template.New("t").Parse(string(idxTPL)))

	err = t.Execute(fd, data)
	if err != nil {
		log.Printf("[FAIL!]write file %s failed: %s", fname, err)

		return
	}

	log.Printf("[SUCC.]write doc to file success: %s", fname)
}

// genDoc generating doc for resource.
func genDoc(dtype, dtypeFolder, fpath, name string, resource *schema.Resource) { //nolint:cyclop,funlen
	data := map[string]string{
		"name":              name,
		"dtype":             strings.ReplaceAll(dtype, "_", ""),
		"dtype_folder":      dtypeFolder,
		"resource":          name[len(providerName)+1:],
		"example":           "",
		"description":       "",
		"description_short": "",
		"import":            "",
	}

	fname := fmt.Sprintf("%s_%s_%s.go", dtype, providerName, data["resource"])
	log.Printf("[START]get description from file: %s\n", fname)

	description, err := getFileDescription(fmt.Sprintf("%s/%s", fpath, fname))
	if err != nil {
		log.Printf("[SKIP!]get description failed, skip: %s", err)

		return
	}

	description = strings.TrimSpace(description)
	if description == "" {
		log.Printf("[SKIP!]description empty, skip: %s\n", fname)

		return
	}

	importPos := strings.Index(description, "\nImport\n")
	if importPos != -1 {
		data["import"] = strings.TrimSpace(description[importPos+8:])
		description = strings.TrimSpace(description[:importPos])
	}

	pos := strings.Index(description, "\nExample Usage\n")
	if pos != -1 {
		data["example"] = strings.TrimSpace(description[pos+15:])
		description = strings.TrimSpace(description[:pos])
	} else {
		log.Printf("[SKIP!]example usage missing, skip: %s\n", fname)

		return
	}

	data["description"] = description

	pos = strings.Index(description, "\n\n")
	if pos != -1 {
		data["description_short"] = strings.TrimSpace(description[:pos])
	} else {
		data["description_short"] = description
	}

	requiredArgs := []string{}
	optionalArgs := []string{}
	attributes := []string{}
	subStruct := []string{}

	keys := make([]string, 0)

	for k := range resource.Schema {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := resource.Schema[k]
		if v.Description == "" {
			continue
		}

		switch {
		case v.Required:
			opt := "Required"
			if v.ForceNew {
				opt += forceNewStr
			}

			requiredArgs = append(
				requiredArgs,
				fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description),
			)

			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		case v.Optional:
			opt := "Optional"
			if v.ForceNew {
				opt += forceNewStr
			}

			optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))

			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		default:
			attrs := getAttributes(0, k, v)
			if len(attrs) > 0 {
				attributes = append(attributes, attrs...)
			}
		}
	}

	sort.Strings(requiredArgs)
	sort.Strings(optionalArgs)
	sort.Strings(attributes)

	requiredArgs = append(requiredArgs, optionalArgs...)
	data["arguments"] = strings.Join(requiredArgs, "\n")

	if len(subStruct) > 0 {
		data["arguments"] += "\n" + strings.Join(subStruct, "\n")
	}

	data["attributes"] = strings.Join(attributes, "\n")

	fname = fmt.Sprintf("%s/%s/%s.md", docRoot, dtypeFolder, data["resource"])

	fd, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644) //nolint:gomnd
	if err != nil {
		log.Printf("[FAIL!]open file %s failed: %s", fname, err)

		return
	}

	defer func() {
		if e := fd.Close(); e != nil {
			log.Printf("[FAIL!]close file %s failed: %s", fname, e)
		}
	}()

	docTPL, err := os.ReadFile(docFile)
	if err != nil {
		log.Printf("[FAIL!]open file %s failed: %s", docFile, err)

		return
	}

	t := template.Must(template.New("t").Parse(string(docTPL)))

	err = t.Execute(fd, data)
	if err != nil {
		log.Printf("[FAIL!]write file %s failed: %s", fname, err)

		return
	}

	log.Printf("[SUCC.]write doc to file success: %s", fname)
}

// getAttributes get attributes from schema.
func getAttributes(step int, k string, v *schema.Schema) []string {
	var attributes []string

	ident := strings.Repeat(" ", step+step)

	if v.Description == "" {
		return attributes
	}

	if v.Computed { //nolint:nestif
		if _, ok := v.Elem.(*schema.Resource); ok {
			var listAttributes []string

			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				attrs := getAttributes(step+1, kk, vv)
				if len(attrs) > 0 {
					listAttributes = append(listAttributes, attrs...)
				}
			}

			slistAttributes := ""

			sort.Strings(listAttributes)

			if len(listAttributes) > 0 {
				slistAttributes = "\n" + strings.Join(listAttributes, "\n")
			}

			attributes = append(
				attributes,
				fmt.Sprintf("%s* `%s` - %s%s", ident, k, v.Description, slistAttributes),
			)
		} else {
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s", ident, k, v.Description))
		}
	}

	return attributes
}

// getFileDescription get description from go file.
func getFileDescription(fname string) (string, error) {
	fset := token.NewFileSet()

	parsedAst, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return parsedAst.Doc.Text(), nil
}

// getSubStruct get sub structure from go file.
//
//nolint:gocognit
func getSubStruct(step int, k string, v *schema.Schema) []string { //nolint:cyclop,funlen
	var subStructs []string

	if v.Description == "" {
		return subStructs
	}

	if v.Type == schema.TypeMap || v.Type == schema.TypeList || v.Type == schema.TypeSet { //nolint:nestif
		if _, ok := v.Elem.(*schema.Resource); ok {
			subStructs = append(
				subStructs,
				fmt.Sprintf("\nThe `%s` object supports the following:\n", k),
			)
			requiredArgs := []string{}
			optionalArgs := []string{}
			attributes := []string{}

			var keys []string
			for kk := range v.Elem.(*schema.Resource).Schema {
				keys = append(keys, kk)
			}

			sort.Strings(keys)

			for _, kk := range keys {
				vv := v.Elem.(*schema.Resource).Schema[kk]
				if vv.Description == "" {
					vv.Description = "************************* Please input Description for Schema ************************* "
				}

				switch {
				case vv.Required:
					opt := "Required"
					if vv.ForceNew {
						opt += forceNewStr
					}

					requiredArgs = append(
						requiredArgs,
						fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description),
					)
				case vv.Optional:
					opt := "Optional"
					if vv.ForceNew {
						opt += forceNewStr
					}

					optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				default:
					attrs := getAttributes(0, kk, vv)
					if len(attrs) > 0 {
						attributes = append(attributes, attrs...)
					}
				}
			}

			sort.Strings(requiredArgs)
			subStructs = append(subStructs, requiredArgs...)

			sort.Strings(optionalArgs)
			subStructs = append(subStructs, optionalArgs...)

			sort.Strings(attributes)
			subStructs = append(subStructs, attributes...)

			for _, kk := range keys {
				vv := v.Elem.(*schema.Resource).Schema[kk]
				subStructs = append(subStructs, getSubStruct(step+1, kk, vv)...)
			}
		}
	}

	return subStructs
}
