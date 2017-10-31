package autoscalr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"math/rand"
	"net/http"
	"time"
)

type AppDef struct {
	AutoScalingGroupName        string   `json:"aws_autoscaling_group_name"`
	AwsRegion                   string   `json:"aws_region"`
	InstanceTypes               []string `json:"instance_types"`
	ScaleMode                   string   `json:"scale_mode"`
	MaxSpotPercentTotal         int      `json:"max_spot_percent_total"`
	MaxSpotPercentOneMarket     int      `json:"max_spot_percent_one_market"`
	TargetSpareCPUPercent       int      `json:"target_spare_cpu_percent"`
	ClusterName                 string   `json:"cluster_name"`
	TargetSpareMemoryPercent    int      `json:"target_spare_memory_percent"`
	QueueName                   string   `json:"queue_name"`
	TargetQueueSize             int      `json:"target_queue_size"`
	MaxMinutesToTargetQueueSize int      `json:"max_minutes_to_target_queue_size"`
	DisplayName                 string   `json:"display_name"`
	DetailedMonitoringEnabled   bool     `json:"detailed_monitoring_enabled"`
	AutoscalrEnabled            bool     `json:"autoscalr_enabled"`
	OsFamily                    string   `json:"os_family"`
	MaxHoursInstanceAge         int      `json:"max_hours_instance_age"`
	TargetCapacity		        int      `json:"target_capacity"`
}

type AutoScalrRequest struct {
	AsrToken    string  `json:"api_key"`
	RequestType string  `json:"request_type"`
	AsrAppDef   *AppDef `json:"autoscalr_app_def"`
}

type AsrApiError struct {
	ErrorMessage    	string  `json:"errorMessage"`
	Code 	 	string  `json:"code"`
}

type AsrApiErrorResponse struct {
	Error    *AsrApiError  `json:"error"`
}

func init() {
	rand.Seed(time.Now().Unix())
}

func resourceAutoScalrAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate,
		Read:   resourceRead,
		Delete: resourceDelete,
		Update: resourceUpdate,

		Schema: map[string]*schema.Schema{
			"aws_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_autoscaling_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_types": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scale_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cpu",
			},
			"max_spot_percent_total": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"max_spot_percent_one_market": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  20,
			},
			"target_spare_cpu_percent": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  20,
			},
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"target_spare_memory_percent": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"queue_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"target_queue_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"max_minutes_to_target_queue_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"detailed_monitoring_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"autoscalr_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"os_family": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Linux/UNIX",
			},
			"max_hours_instance_age": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"target_capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func makeApiCall(d *schema.ResourceData, meta interface{}, asrReq *AutoScalrRequest, resId string) (int, *AppDef, error) {
	config := meta.(*Config)
	url := config.apiUrl
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(asrReq)
	app := new(AppDef)
	resp, err := client.Post(url, "application/json", postBody)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// make 2 copies of response, one for error decoding and one for good response
			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)
			errBuf := bytes.NewBuffer(respBuf.Bytes())
			// Check for error response json
			jsonErr := new(AsrApiErrorResponse)
			json.NewDecoder(errBuf).Decode(jsonErr)
			if jsonErr.Error != nil && jsonErr.Error.ErrorMessage != ""  {
				// error response
				err = errors.New(fmt.Sprintf("Error response: %s", jsonErr.Error.ErrorMessage))
			} else {
				// looks like good response
				json.NewDecoder(respBuf).Decode(app)
				d.SetId(resId)
			}
			return resp.StatusCode, app, err
		} else {
			err = errors.New(fmt.Sprintf("AutoScalr API returned: %d", resp.Status))
			return resp.StatusCode, app, err
		}
	} else {
		//log.Println("Error: %s", err.Error())
		return 500, app, err
	}
}

func resourceCreate(d *schema.ResourceData, meta interface{}) error {
	awsRegion := d.Get("aws_region").(string)
	autoScalingGroupName := d.Get("aws_autoscaling_group_name").(string)
	instanceTypesTemp := d.Get("instance_types").([]interface{})
	instanceTypes := make([]string, len(instanceTypesTemp))
	for i, v := range instanceTypesTemp {
		instanceTypes[i] = v.(string)
	}
	scaleMode := d.Get("scale_mode").(string)
	maxSpotPercentTotal := d.Get("max_spot_percent_total").(int)
	maxSpotPercentOneMarket := d.Get("max_spot_percent_one_market").(int)
	clusterName := d.Get("cluster_name").(string)
	targetSpareCpuPercent := d.Get("target_spare_cpu_percent").(int)
	targetSpareMemoryPercent := d.Get("target_spare_memory_percent").(int)
	queueName := d.Get("queue_name").(string)
	targetQueueSize := d.Get("target_queue_size").(int)
	maxMinutesToTargetQueueSize := d.Get("max_minutes_to_target_queue_size").(int)
	displayName := d.Get("display_name").(string)
	detailedMonitoringEnabled := d.Get("detailed_monitoring_enabled").(bool)
	autoscalrEnabled := d.Get("autoscalr_enabled").(bool)
	osFamily := d.Get("os_family").(string)
	maxHoursInstanceAge := d.Get("max_hours_instance_age").(int)
	targetCapacity := d.Get("target_capacity").(int)

	config := meta.(*Config)

	body := &AutoScalrRequest{
		AsrToken:    config.AccessKey,
		RequestType: "Create",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        autoScalingGroupName,
			AwsRegion:                   awsRegion,
			InstanceTypes:               instanceTypes,
			ScaleMode:                   scaleMode,
			MaxSpotPercentTotal:         maxSpotPercentTotal,
			MaxSpotPercentOneMarket:     maxSpotPercentOneMarket,
			ClusterName:                 clusterName,
			TargetSpareCPUPercent:       targetSpareCpuPercent,
			TargetSpareMemoryPercent:    targetSpareMemoryPercent,
			QueueName:                   queueName,
			TargetQueueSize:             targetQueueSize,
			MaxMinutesToTargetQueueSize: maxMinutesToTargetQueueSize,
			DisplayName:                 displayName,
			DetailedMonitoringEnabled:   detailedMonitoringEnabled,
			AutoscalrEnabled:            autoscalrEnabled,
			OsFamily:                    osFamily,
			MaxHoursInstanceAge:         maxHoursInstanceAge,
			TargetCapacity:         	 targetCapacity,
		},
	}
	resId := fmt.Sprintf("%s:%s", autoScalingGroupName, awsRegion)

	respCode, _, err := makeApiCall(d, meta, body, resId)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}

func resourceUpdate(d *schema.ResourceData, meta interface{}) error {
	awsRegion := d.Get("aws_region").(string)
	autoScalingGroupName := d.Get("aws_autoscaling_group_name").(string)
	instanceTypesTemp := d.Get("instance_types").([]interface{})
	instanceTypes := make([]string, len(instanceTypesTemp))
	for i, v := range instanceTypesTemp {
		instanceTypes[i] = v.(string)
	}
	scaleMode := d.Get("scale_mode").(string)
	maxSpotPercentTotal := d.Get("max_spot_percent_total").(int)
	maxSpotPercentOneMarket := d.Get("max_spot_percent_one_market").(int)
	clusterName := d.Get("cluster_name").(string)
	targetSpareCpuPercent := d.Get("target_spare_cpu_percent").(int)
	targetSpareMemoryPercent := d.Get("target_spare_memory_percent").(int)
	queueName := d.Get("queue_name").(string)
	targetQueueSize := d.Get("target_queue_size").(int)
	maxMinutesToTargetQueueSize := d.Get("max_minutes_to_target_queue_size").(int)
	displayName := d.Get("display_name").(string)
	detailedMonitoringEnabled := d.Get("detailed_monitoring_enabled").(bool)
	autoscalrEnabled := d.Get("autoscalr_enabled").(bool)
	osFamily := d.Get("os_family").(string)
	maxHoursInstanceAge := d.Get("max_hours_instance_age").(int)
	targetCapacity := d.Get("target_capacity").(int)

	config := meta.(*Config)

	body := &AutoScalrRequest{
		AsrToken:    config.AccessKey,
		RequestType: "Update",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        autoScalingGroupName,
			AwsRegion:                   awsRegion,
			InstanceTypes:               instanceTypes,
			ScaleMode:                   scaleMode,
			MaxSpotPercentTotal:         maxSpotPercentTotal,
			MaxSpotPercentOneMarket:     maxSpotPercentOneMarket,
			ClusterName:                 clusterName,
			TargetSpareCPUPercent:       targetSpareCpuPercent,
			TargetSpareMemoryPercent:    targetSpareMemoryPercent,
			QueueName:                   queueName,
			TargetQueueSize:             targetQueueSize,
			MaxMinutesToTargetQueueSize: maxMinutesToTargetQueueSize,
			DisplayName:                 displayName,
			DetailedMonitoringEnabled:   detailedMonitoringEnabled,
			AutoscalrEnabled:            autoscalrEnabled,
			OsFamily:                    osFamily,
			MaxHoursInstanceAge:         maxHoursInstanceAge,
			TargetCapacity:         	 targetCapacity,
		},
	}
	resId := fmt.Sprintf("%s:%s", autoScalingGroupName, awsRegion)
	respCode, _, err := makeApiCall(d, meta, body, resId)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	awsRegion := d.Get("aws_region").(string)
	autoScalingGroupName := d.Get("aws_autoscaling_group_name").(string)
	body := &AutoScalrRequest{
		AsrToken:    config.AccessKey,
		RequestType: "Get",
		AsrAppDef: &AppDef{
			AutoScalingGroupName: autoScalingGroupName,
			AwsRegion:            awsRegion,
		},
	}
	resId := fmt.Sprintf("%s:%s", autoScalingGroupName, awsRegion)
	respCode, app, err := makeApiCall(d, meta, body, resId)
	if respCode > 0 {
		//log.Println("respBody:" + respBody)
		if respCode == 200 {
			if app.AutoScalingGroupName == autoScalingGroupName {
				// resource still exists, update values to returned values
				d.Set("instance_types", app.InstanceTypes)
				d.Set("scale_mode", app.ScaleMode)
				d.Set("max_spot_percent_total", app.MaxSpotPercentTotal)
				d.Set("max_spot_percent_one_market", app.MaxSpotPercentOneMarket)
				d.Set("target_spare_cpu_percent", app.TargetSpareCPUPercent)
				d.Set("cluster_name", app.ClusterName)
				d.Set("target_spare_memory_percent", app.TargetSpareMemoryPercent)
				d.Set("queue_name", app.QueueName)
				d.Set("target_queue_size", app.TargetQueueSize)
				d.Set("max_minutes_to_target_queue_size", app.MaxMinutesToTargetQueueSize)
				d.Set("display_name", app.DisplayName)
				d.Set("detailed_monitoring_enabled", app.DetailedMonitoringEnabled)
				d.Set("autoscalr_enabled", app.AutoscalrEnabled)
				d.Set("os_family", app.OsFamily)
				d.Set("max_hours_instance_age", app.MaxHoursInstanceAge)
				d.Set("target_capacity", app.TargetCapacity)
			} else {
				// resource must have been deleted out of band
				// Set id to tell terraform
				d.SetId("")
			}
		} else {
			err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
		}
	}
	return err
}

func resourceDelete(d *schema.ResourceData, meta interface{}) error {
	awsRegion := d.Get("aws_region").(string)
	autoScalingGroupName := d.Get("aws_autoscaling_group_name").(string)
	instanceTypesTemp := d.Get("instance_types").([]interface{})
	instanceTypes := make([]string, len(instanceTypesTemp))
	for i, v := range instanceTypesTemp {
		instanceTypes[i] = v.(string)
	}
	scaleMode := d.Get("scale_mode").(string)
	maxSpotPercentTotal := d.Get("max_spot_percent_total").(int)
	maxSpotPercentOneMarket := d.Get("max_spot_percent_one_market").(int)
	clusterName := d.Get("cluster_name").(string)
	targetSpareCpuPercent := d.Get("target_spare_cpu_percent").(int)
	targetSpareMemoryPercent := d.Get("target_spare_memory_percent").(int)
	queueName := d.Get("queue_name").(string)
	targetQueueSize := d.Get("target_queue_size").(int)
	maxMinutesToTargetQueueSize := d.Get("max_minutes_to_target_queue_size").(int)
	displayName := d.Get("display_name").(string)
	detailedMonitoringEnabled := d.Get("detailed_monitoring_enabled").(bool)
	autoscalrEnabled := d.Get("autoscalr_enabled").(bool)
	osFamily := d.Get("os_family").(string)
	maxHoursInstanceAge := d.Get("max_hours_instance_age").(int)
	targetCapacity := d.Get("target_capacity").(int)

	config := meta.(*Config)

	body := &AutoScalrRequest{
		AsrToken:    config.AccessKey,
		RequestType: "Delete",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        autoScalingGroupName,
			AwsRegion:                   awsRegion,
			InstanceTypes:               instanceTypes,
			ScaleMode:                   scaleMode,
			MaxSpotPercentTotal:         maxSpotPercentTotal,
			MaxSpotPercentOneMarket:     maxSpotPercentOneMarket,
			ClusterName:                 clusterName,
			TargetSpareCPUPercent:       targetSpareCpuPercent,
			TargetSpareMemoryPercent:    targetSpareMemoryPercent,
			QueueName:                   queueName,
			TargetQueueSize:             targetQueueSize,
			MaxMinutesToTargetQueueSize: maxMinutesToTargetQueueSize,
			DisplayName:                 displayName,
			DetailedMonitoringEnabled:   detailedMonitoringEnabled,
			AutoscalrEnabled:            autoscalrEnabled,
			OsFamily:                    osFamily,
			MaxHoursInstanceAge:         maxHoursInstanceAge,
			TargetCapacity:         	 targetCapacity,
		},
	}
	resId := fmt.Sprintf("%s:%s", autoScalingGroupName, awsRegion)

	respCode, _, err := makeApiCall(d, meta, body, resId)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}
