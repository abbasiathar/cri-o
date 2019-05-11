package server_test

import (
	"context"

	"github.com/cri-o/cri-o/oci"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// The actual test suite
var _ = t.Describe("ListPodSandbox", func() {
	// Prepare the sut
	BeforeEach(func() {
		beforeEach()
		setupSUT()
	})

	AfterEach(afterEach)

	t.Describe("ListPodSandbox", func() {
		It("should succeed", func() {
			// Given
			sut.AddSandbox(testSandbox)
			testContainer.SetState(&oci.ContainerState{
				State: specs.State{Status: oci.ContainerStateRunning},
			})
			testSandbox.SetCreated()
			Expect(testSandbox.SetInfraContainer(testContainer)).To(BeNil())

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeEquivalentTo(1))
		})

		It("should succeed without infra container", func() {
			// Given
			sut.AddSandbox(testSandbox)
			testSandbox.SetCreated()

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeZero())
		})

		It("should skip not created sandboxes", func() {
			// Given
			sut.AddSandbox(testSandbox)
			Expect(testSandbox.SetInfraContainer(testContainer)).To(BeNil())

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeZero())
		})

		It("should succeed with filter", func() {
			// Given
			mockDirs(testManifest)
			Expect(sut.LoadSandbox(sandboxID)).To(BeNil())

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{Filter: &pb.PodSandboxFilter{
					Id: sandboxID,
				}})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeEquivalentTo(1))
		})

		It("should succeed with filter for state", func() {
			// Given
			mockDirs(testManifest)
			Expect(sut.LoadSandbox(sandboxID)).To(BeNil())

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{Filter: &pb.PodSandboxFilter{
					Id: sandboxID,
					State: &pb.PodSandboxStateValue{
						State: pb.PodSandboxState_SANDBOX_READY},
				}})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeZero())
		})

		It("should succeed with filter for label", func() {
			// Given
			mockDirs(testManifest)
			Expect(sut.LoadSandbox(sandboxID)).To(BeNil())

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{Filter: &pb.PodSandboxFilter{
					Id:            sandboxID,
					LabelSelector: map[string]string{"label": "value"},
				}})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeZero())
		})

		It("should succeed with filter but when not finding id", func() {
			// Given
			sut.AddSandbox(testSandbox)

			// When
			response, err := sut.ListPodSandbox(context.Background(),
				&pb.ListPodSandboxRequest{Filter: &pb.PodSandboxFilter{
					Id: sandboxID,
				}})

			// Then
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(len(response.Items)).To(BeZero())
		})
	})
})