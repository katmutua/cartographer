// Copyright 2021 VMware
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

package v1alpha1_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
)

var _ = Describe("ClusterRunTemplate", func() {
	Describe("Webhook Validation", func() {
		var (
			template *v1alpha1.ClusterRunTemplate
		)

		BeforeEach(func() {
			template = &v1alpha1.ClusterRunTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-template",
					Namespace: "default",
				},
			}
		})

		Describe("#Create", func() {
			Context("template is well formed", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(&ArbitraryObject{
						TypeMeta: metav1.TypeMeta{
							Kind:       "some-kind",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "some-name",
						},
						Spec: ArbitrarySpec{
							SomeKey: "some-val",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("succeeds", func() {
					Expect(template.ValidateCreate()).To(Succeed())
				})
			})

			Context("template sets object namespace", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(&ArbitraryObject{
						TypeMeta: metav1.TypeMeta{
							Kind:       "some-kind",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "some-name",
							Namespace: "some-namespace",
						},
						Spec: ArbitrarySpec{
							SomeKey: "some-val",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("returns an error", func() {
					Expect(template.ValidateCreate()).
						To(MatchError("invalid template: template should not set metadata.namespace on the child object"))
				})
			})

			Context("templated object does not have a spec", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(v1.ConfigMap{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ConfigMap",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "another-name",
						},
						Data: map[string]string{
							"greeting":   "hi",
							"salutation": "bye",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("returns an error", func() {
					Expect(template.ValidateCreate()).
						To(MatchError(ContainSubstring("invalid template: object must have a spec; templated object:")))
				})
			})
		})

		Describe("#Update", func() {
			Context("template is well formed", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(&ArbitraryObject{
						TypeMeta: metav1.TypeMeta{
							Kind:       "some-kind",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "some-name",
						},
						Spec: ArbitrarySpec{
							SomeKey: "some-val",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("succeeds", func() {
					Expect(template.ValidateUpdate(nil)).To(Succeed())
				})
			})

			Context("template sets object namespace", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(&ArbitraryObject{
						TypeMeta: metav1.TypeMeta{
							Kind:       "some-kind",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "some-name",
							Namespace: "some-namespace",
						},
						Spec: ArbitrarySpec{
							SomeKey: "some-val",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("returns an error", func() {
					Expect(template.ValidateUpdate(nil)).
						To(MatchError("invalid template: template should not set metadata.namespace on the child object"))
				})
			})

			Context("templated object does not have a spec", func() {
				BeforeEach(func() {
					raw, err := json.Marshal(v1.ConfigMap{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ConfigMap",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "another-name",
						},
						Data: map[string]string{
							"greeting":   "hi",
							"salutation": "bye",
						},
					})
					Expect(err).NotTo(HaveOccurred())
					template.Spec.Template = runtime.RawExtension{Raw: raw}
				})

				It("returns an error", func() {
					Expect(template.ValidateCreate()).
						To(MatchError(ContainSubstring("invalid template: object must have a spec; templated object:")))
				})
			})
		})

		Context("#Delete", func() {
			Context("Any template", func() {
				var anyTemplate *v1alpha1.ClusterRunTemplate
				It("always succeeds", func() {
					Expect(anyTemplate.ValidateDelete()).NotTo(HaveOccurred())
				})
			})
		})
	})
})
