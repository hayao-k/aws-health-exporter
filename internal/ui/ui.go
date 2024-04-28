package ui

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/health"
	"github.com/aws/aws-sdk-go-v2/service/health/types"
	"github.com/manifoldco/promptui"
)

func PromptEventFilters() (string, string) {
	prompt := promptui.Prompt{
		Label:     "Enter service (or leave blank for all services)",
		Default:   "",
		AllowEdit: true,
	}
	service, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return "", ""
	}

	prompt = promptui.Prompt{
		Label:     "Enter status code (or leave blank for all statuses)",
		Default:   "",
		AllowEdit: true,
	}
	status, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return "", ""
	}

	return service, status
}

func PrepareEventChoicesAndMap(
	eventsResp *health.DescribeEventsForOrganizationOutput,
) (eventChoices []string, eventsMap map[string]types.OrganizationEvent) {
	eventChoices = make([]string, len(eventsResp.Events))
	eventsMap = make(map[string]types.OrganizationEvent, len(eventsResp.Events))
	for i, e := range eventsResp.Events {
		startTime := "N/A"
		if e.StartTime != nil {
			startTime = e.StartTime.Format("2006-01-02 15:04:05")
		}
		choice := fmt.Sprintf(
			"%s - %s (%s, %s)",
			aws.ToString(e.Service),
			aws.ToString(e.EventTypeCode),
			aws.ToString(e.Region),
			startTime,
		)
		eventChoices[i] = choice
		eventsMap[choice] = e
	}
	return
}

func PromptEventSelection(
	eventChoices []string,
	eventsMap map[string]types.OrganizationEvent,
) (types.OrganizationEvent, error) {
	prompt := promptui.Select{
		Label: "Select an event",
		Items: eventChoices,
	}
	_, selectedEventStr, err := prompt.Run()
	if err != nil {
		return types.OrganizationEvent{},
			fmt.Errorf("Event selection canceled: %w", err)
	}
	selectedEvent := eventsMap[selectedEventStr]
	return selectedEvent, nil
}
