package sample_function_test

import (
	"path"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/pivotal-cf/pfs-system-test/util"
)

var _ = Describe("SampleFunctionTest", func() {
	Describe("SampleFunctionTest", func() {

		Context("deploy sample greeter java function", func() {

			It("builds and deploys", func() {
				functionDir := path.Join(util.TEST_CONFIG.BaseDir, "samples", "java", "greeter")

				functionName := util.RandStringShort()
				inputTopicName := util.RandStringShort()
				outputTopicName := util.RandStringShort()
				imageName := util.TEST_CONFIG.DockerOrg + "/" + functionName
				workloadFileSource := path.Join(functionDir, "greeter.yaml")
				workloadFileTarget := path.Join(functionDir, functionName+".yaml")

				util.MvnCleanPackage(functionDir)
				util.CopyAndReplace(path.Join(functionDir, "Dockerfile"), path.Join(functionDir, "Dockerfile"), "java-function-invoker:.*", "java-function-invoker:"+util.TEST_CONFIG.JavaInvokerVersion)

				util.DockerBuild(functionDir, imageName)
				util.DockerPush(imageName)
				util.CopyAndReplace(workloadFileSource, workloadFileTarget, "name: greeter", "name: "+functionName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "name: names", "name: "+inputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "name: greetings", "name: "+outputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "input: names", "input: "+inputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "output: greetings", "output: "+outputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "image: projectriff/greeter:.*", "image: "+imageName)

				util.KubectlApply(workloadFileTarget, util.TEST_CONFIG.Namespace)
				util.SendMessageToGateway(inputTopicName, "World")

				outputMessage := util.KubectlFromKafkaPod(outputTopicName)
				gomega.Expect(outputMessage).To(gomega.MatchRegexp(`(?s:.*Hello World.*)`))

				util.KubectlDelete(workloadFileTarget, util.TEST_CONFIG.Namespace)
				util.DeleteFile(workloadFileTarget)
			})
		})

		Context("deploy sample square node function", func() {

			It("builds and deploys", func() {
				functionDir := path.Join(util.TEST_CONFIG.BaseDir, "samples", "node", "square")

				functionName := util.RandStringShort()
				inputTopicName := util.RandStringShort()
				imageName := util.TEST_CONFIG.DockerOrg + "/" + functionName
				workloadFileSource := path.Join(functionDir, "square.yaml")
				workloadFileTarget := path.Join(functionDir, functionName+".yaml")

				util.CopyAndReplace(path.Join(functionDir, "Dockerfile"), path.Join(functionDir, "Dockerfile"), "node-function-invoker:.*", "node-function-invoker:"+util.TEST_CONFIG.NodeInvokerVersion)

				util.DockerBuild(functionDir, imageName)
				util.DockerPush(imageName)
				util.CopyAndReplace(workloadFileSource, workloadFileTarget, "name: square", "name: "+functionName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "name: numbers", "name: "+inputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "input: numbers", "input: "+inputTopicName)
				util.CopyAndReplace(workloadFileTarget, workloadFileTarget, "image: projectriff/square:.*", "image: "+imageName)

				util.KubectlApply(workloadFileTarget, util.TEST_CONFIG.Namespace)
				reply := util.SendRequestToGateway(inputTopicName, "12")

				gomega.Expect(reply).To(gomega.MatchRegexp(`(?s:.*144.*)`))

				util.KubectlDelete(workloadFileTarget, util.TEST_CONFIG.Namespace)
				util.DeleteFile(workloadFileTarget)
			})
		})
	})
})
