// +build ignore

package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"sort"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"
)

const awsTmpl = `
// Code generated by generate.go; DO NOT EDIT.

package constants

var Zones = []*Zone{
{{ range $zone := . -}}
	&Zone{"{{$zone.Region}}", []string{
		{{ range $az := $zone.AvailabilityZones -}}
			"{{$az}}",
		{{end -}}
	}},
{{ end }}
}

type Zone struct {
	Region            string
	AvailabilityZones []string
}
`

type Zone struct {
	Region            string
	AvailabilityZones []string
}

func main() {
	err := createAwsConstants()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func createAwsConstants() error {
	zones, err := loadZones()
	if err != nil {
		return err
	}

	err = writeAwsConstants("./constants/aws.go", zones)
	if err != nil {
		return err
	}

	return nil
}

func loadZones() ([]*Zone, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load aws config")
	}
	svc := ec2.NewFromConfig(cfg)

	regionInput := ec2.DescribeRegionsInput{}
	res, err := svc.DescribeRegions(context.TODO(), &regionInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load aws regions")
	}

	var zones []*Zone
	for _, r := range res.Regions {
		zones = append(zones, &Zone{Region: *r.RegionName})
	}

	azInput := &ec2.DescribeAvailabilityZonesInput{}
	for _, z := range zones {
		fmt.Printf("loading availability zones for %s...\n", z.Region)

		result, err := svc.DescribeAvailabilityZones(context.TODO(), azInput, func(options *ec2.Options) {
			options.Region = z.Region
		})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to load availability zones for %s", z.Region))
		}
		var azs []string
		for _, az := range result.AvailabilityZones {
			azs = append(azs, *az.ZoneName)
		}

		sort.Strings(azs)
		z.AvailabilityZones = azs
	}

	return sortZones(zones), nil
}

func sortZones(zones []*Zone) []*Zone {
	var regions []string
	zoneMap := map[string]*Zone{}
	for _, zone := range zones {
		regions = append(regions, zone.Region)
		zoneMap[zone.Region] = zone
	}

	sort.Strings(regions)
	var sortedZones []*Zone
	for _, r := range regions {
		sortedZones = append(sortedZones, zoneMap[r])
	}

	return sortedZones
}

func writeAwsConstants(filename string, zones []*Zone) error {
	tmpl, err := template.New("aws").Parse(awsTmpl)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse template. text=\n %s", awsTmpl))
	}

	writer := bytes.Buffer{}
	if err = tmpl.Execute(&writer, zones); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to apply template\nzones: %v\ntext=\n %s", zones, awsTmpl))
	}

	formatted, err := format.Source(writer.Bytes())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to format text=\n %s", writer.String()))
	}

	err = ioutil.WriteFile(filename, formatted, 0776)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write generated text to %s", filename))
	}

	return nil
}
