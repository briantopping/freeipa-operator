# FreeIPA Operator
[![GitHub license][license-badge]](LICENSE)

An integrated Identity and Authentication solution for Linux/UNIX networked environments. IPA stands for "Identity, Policy, Audit".
Comparatively, it provides a solution similar to Microsoft Active Directory for *NIX environments.

At it's core, FreeIPA uses the 389 LDAP server for storage of all client information. 389 has solid replication capabilities,
but it's data sets depend on IP addresses that stay coupled to specific instance data. This creates challenges in an ephemeral 
container environment. 

The FreeIPA operator works to bridge these challenges so users can focus on the value proposition of FreeIPA instead
of spending weeks getting it settled and reliable in a continerized environment. 

## Tooling and Dependencies

The operator is based on Kubebuilder and Bazel. The container itself comes from the FreeIPA Container project. 

Dependencies for the project are provided by `dep`. `dep ensure` will bring in the vendor directory based on the `Gopkg.toml`
file at the root of the project. At some point, this should not be necessary for non-developer builds as Bazel will be 
configured in `WORKSPACE` to pull the dependencies. Changes to the dependencies will use a combination of `dep ensure`
and Gazelle to update the Bazel builds.

## Building

Until the Bazel build is finished, use the `Makefile` for tasks like generating files in `config` and `pkg/apis`:
* `make generate` - Will create the CRDs and deepcopy routines. This is good if you want to run from your IDE debugger.
* `make install` - Install the generated CRD to the current cluster context.
* `make run` - This will run the controller against the cluster currently configured in your `~/.kube/config`. In this case,
debugging is up to you (such as connecting Delve to the running process).

In general, once files are generated, the `cmd/manager/main.go` can be run from your IDE for debugging.

Please take a closer look at the `Makefile` for other options.

## Contributing

Please feel free to engage with issues and PRs. It's imagined this project will reach basic maturity pretty quickly, but there's
a lot of benefits that the operator pattern could bring to FreeIPA, such as managing backups, upgrades and monitoring with 
Prometheus.   

## Status

This project is in development phase. We'll keep this page updated as we get usable features.

## Code of conduct

This project is for everyone. We ask that our users and contributors take a few
minutes to review our [code of conduct][coc].


## License

FreeIPA Operator is copyright 2018 The FreeIPA Operator Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

<!-- refs -->
[coc]: https://github.com/linkerd/linkerd/wiki/Linkerd-code-of-conduct
[license-badge]: https://img.shields.io/github/license/linkerd/linkerd.svg
