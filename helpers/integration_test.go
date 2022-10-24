package helpers_test

import (
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-appstudio/integration-service/helpers"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Helpers for integration", func() {

	var pipelineRun *tektonv1beta1.PipelineRun

	BeforeEach(func() {
		pipelineRun = &tektonv1beta1.PipelineRun{
			ObjectMeta: v1.ObjectMeta{
				Name:      "example-pipeline-run",
				Namespace: "default",
			},
			Status: tektonv1beta1.PipelineRunStatus{
				PipelineRunStatusFields: tektonv1beta1.PipelineRunStatusFields{
					TaskRuns: map[string]*tektonv1beta1.PipelineRunTaskRunStatus{
						"task-1": {
							PipelineTaskName: "pr-1-task-1",
							Status: &tektonv1beta1.TaskRunStatus{
								TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
									TaskRunResults: []tektonv1beta1.TaskRunResult{{
										Name:  "HACBS_TEST_OUTPUT",
										Value: "{\"result\": \"SUCCESS\", \"timestamp\": \"1640995200\", \"failures\": 0, \"successes\": 10}",
									}},
								},
							},
						},
						"task-2": {
							PipelineTaskName: "pr-1-task-2",
							Status: &tektonv1beta1.TaskRunStatus{
								TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
									TaskRunResults: []tektonv1beta1.TaskRunResult{{
										Name:  "HACBS_TEST_OUTPUT",
										Value: "{\"result\": \"SKIPPED\", \"timestamp\": \"1640995200\", \"failures\": 0, \"successes\": 1}",
									}},
								},
							},
						},
					},
				},
			},
		}
	})

	It("can calculate the integration PipelineRun outcome", func() {
		result, err := helpers.CalculateIntegrationPipelineRunOutcome(logr.Discard(), pipelineRun)
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		pipelineRun.Status.PipelineRunStatusFields.TaskRuns["task-2"].Status.TaskRunStatusFields = tektonv1beta1.TaskRunStatusFields{
			TaskRunResults: []tektonv1beta1.TaskRunResult{{
				Name:  "HACBS_TEST_OUTPUT",
				Value: "{\"result\": \"FAILED\", \"timestamp\": \"1640995200\", \"failures\": 5, \"successes\": 0}",
			}},
		}

		result, err = helpers.CalculateIntegrationPipelineRunOutcome(logr.Discard(), pipelineRun)
		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())
	})

	It("can handle malformed HACBS_TEST_OUTPUT result", func() {
		pipelineRun.Status.PipelineRunStatusFields.TaskRuns["task-2"].Status.TaskRunStatusFields = tektonv1beta1.TaskRunStatusFields{
			TaskRunResults: []tektonv1beta1.TaskRunResult{{
				Name:  "HACBS_TEST_OUTPUT",
				Value: "invalid json",
			}},
		}

		result, err := helpers.CalculateIntegrationPipelineRunOutcome(logr.Discard(), pipelineRun)
		Expect(err).ToNot(BeNil())
		Expect(result).To(BeFalse())
	})
})
