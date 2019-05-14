load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.18.3/rules_go-0.18.3.tar.gz"],
    sha256 = "86ae934bd4c43b99893fc64be9d9fc684b81461581df7ea8fc291c816f5ee8c5",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/archive/ddff739eca1ac0d5cc3530fc6af2ccfe19becc78.tar.gz"],
    strip_prefix = "bazel-gazelle-ddff739eca1ac0d5cc3530fc6af2ccfe19becc78",
    sha256 = "59e1fd653eeb4deebdafd73c3dbd55370a7d63a60cdbe51472e33c36996f3dd9",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load(
    "@io_bazel_rules_go//proto:def.bzl",
    "proto_register_toolchains",
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "io_k8s_api",
    build_file_generation = "on",
    build_file_proto_mode = "legacy",
    commit = "b503174bad5991eb66f18247f52e41c3258f6348",
    importpath = "k8s.io/api",
)

go_repository(
    name = "io_k8s_apiextensions_apiserver",
    build_file_generation = "on",
    build_file_proto_mode = "legacy",
    commit = "0cd23ebeb6882bd1cdc2cb15fc7b2d72e8a86a5b",
    importpath = "k8s.io/apiextensions-apiserver",
    build_extra_args = ["-exclude=vendor"],
)

go_repository(
    name = "io_k8s_apimachinery",
    build_file_generation = "on",
    build_file_proto_mode = "legacy",
    commit = "eddba98df674a16931d2d4ba75edc3a389bf633a",
    importpath = "k8s.io/apimachinery",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "4b2b341e8d7715fae06375aa633dbb6e91b3fb46",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    commit = "8991bc29aa16c548c550c7ff78260e27b9ab7c73",
    importpath = "github.com/davecgh/go-spew",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    commit = "103c9496ad8f7e687b8291b56750190012091a96",
    importpath = "github.com/emicklei/go-restful",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "com_github_go_logr_logr",
    commit = "9fb12b3b21c5415d16ac18dc5cd42c1cfdd40c4e",
    importpath = "github.com/go-logr/logr",
)

go_repository(
    name = "com_github_go_logr_zapr",
    commit = "03f06a783fbb7dfaf3f629c7825480e43a7105e6",
    importpath = "github.com/go-logr/zapr",
)

go_repository(
    name = "com_github_gobuffalo_envy",
    commit = "043cb4b8af871b49563291e32c66bb84378a60ac",
    importpath = "github.com/gobuffalo/envy",
)

go_repository(
    name = "com_github_gogo_protobuf",
    commit = "ba06b47c162d49f2af050fb4c75bcbc86a159d5c",
    importpath = "github.com/gogo/protobuf",
)

go_repository(
    name = "com_github_golang_glog",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    importpath = "github.com/golang/glog",
)

go_repository(
    name = "com_github_golang_groupcache",
    commit = "5b532d6fd5efaf7fa130d4e859a2fde0fc3a9e1b",
    importpath = "github.com/golang/groupcache",
)

go_repository(
    name = "com_github_golang_protobuf",
    commit = "b5d812f8a3706043e23a9cd5babf2e5423744d30",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_google_btree",
    commit = "4030bb1f1f0c35b30ca7009e9ebd06849dd45306",
    importpath = "github.com/google/btree",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "f140a6486e521aad38f5917de355cbf147cc0496",
    importpath = "github.com/google/gofuzz",
)

go_repository(
    name = "com_github_google_uuid",
    commit = "0cd6bf5da1e1c83f8b45653022c74f71af0538a4",
    importpath = "github.com/google/uuid",
)

go_repository(
    name = "com_github_googleapis_gnostic",
    commit = "7c663266750e7d82587642f65e60bc4083f1f84e",
    importpath = "github.com/googleapis/gnostic",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    commit = "3befbb6ad0cc97d4c25d851e9528915809e1a22f",
    importpath = "github.com/gregjones/httpcache",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    commit = "7087cb70de9f7a8bc0a10c375cb0d2280a8edf9c",
    importpath = "github.com/hashicorp/golang-lru",
)

go_repository(
    name = "com_github_hpcloud_tail",
    commit = "a30252cb686a21eb2d0b98132633053ec2f7f1e5",
    importpath = "github.com/hpcloud/tail",
)

go_repository(
    name = "com_github_imdario_mergo",
    commit = "7c29201646fa3de8506f701213473dd407f19646",
    importpath = "github.com/imdario/mergo",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    commit = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75",
    importpath = "github.com/inconshreveable/mousetrap",
)

go_repository(
    name = "com_github_joho_godotenv",
    commit = "23d116af351c84513e1946b527c88823e476be13",
    importpath = "github.com/joho/godotenv",
)

go_repository(
    name = "com_github_json_iterator_go",
    commit = "0ff49de124c6f76f8494e194af75bde0f1a49a29",
    importpath = "github.com/json-iterator/go",
)

go_repository(
    name = "com_github_markbates_inflect",
    commit = "24b83195037b3bc61fcda2d28b7b0518bce293b6",
    importpath = "github.com/markbates/inflect",
)

go_repository(
    name = "com_github_mattbaird_jsonpatch",
    commit = "81af80346b1a01caae0cbc27fd3c1ba5b11e189f",
    importpath = "github.com/mattbaird/jsonpatch",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_modern_go_concurrent",
    commit = "bacd9c7ef1dd9b15be4a9909b8ac7a4e313eec94",
    importpath = "github.com/modern-go/concurrent",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    commit = "4b7aa43c6742a2c18fdef89dd197aaae7dac7ccd",
    importpath = "github.com/modern-go/reflect2",
)

go_repository(
    name = "com_github_onsi_ginkgo",
    commit = "eea6ad008b96acdaa524f5b409513bf062b500ad",
    importpath = "github.com/onsi/ginkgo",
)

go_repository(
    name = "com_github_onsi_gomega",
    commit = "90e289841c1ed79b7a598a7cd9959750cb5e89e2",
    importpath = "github.com/onsi/gomega",
)

go_repository(
    name = "com_github_pborman_uuid",
    commit = "adf5a7427709b9deb95d29d3fa8a2bf9cfd388f1",
    importpath = "github.com/pborman/uuid",
)

go_repository(
    name = "com_github_petar_gollrb",
    commit = "53be0d36a84c2a886ca057d34b6aa4468df9ccb4",
    importpath = "github.com/petar/GoLLRB",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    commit = "0be1b92a6df0e4f5cb0a5d15fb7f643d0ad93ce6",
    importpath = "github.com/peterbourgon/diskv",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "ba968bfe8b2f7e042a574c888954fccecfa385b4",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "505eaef017263e299324067d40ca2c48f6a2cf50",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "fd36f4220a901265f90734c3183c5f0c91daa0b8",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "1ba88736f028e37bc17328369e94a537ae9e0234",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "5867b95ac084bbfee6ea16595c4e05ab009021da",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_rogpeppe_go_internal",
    commit = "438578804ca6f31be148c27683afc419ce47c06e",
    importpath = "github.com/rogpeppe/go-internal",
)

go_repository(
    name = "com_github_spf13_afero",
    commit = "588a75ec4f32903aa5e39a2619ba6a4631e28424",
    importpath = "github.com/spf13/afero",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "ef82de70bb3f60c65fb8eebacbb2d122ef517385",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "298182f68c66c05229eb03ac171abe6e309ee79a",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "8c41231e01b2085512d98153bcffb847ff9b4b9f",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "in_gopkg_fsnotify_v1",
    commit = "c2828203cd70a50dcccfb2761f8b1f8ceef9a8e9",
    importpath = "gopkg.in/fsnotify.v1",
    remote = "https://github.com/fsnotify/fsnotify.git",
    vcs = "git",
)

go_repository(
    name = "in_gopkg_inf_v0",
    commit = "d2d2541c53f18d2a059457998ce2876cc8e67cbf",
    importpath = "gopkg.in/inf.v0",
)

go_repository(
    name = "in_gopkg_tomb_v1",
    commit = "dd632973f1e7218eb1089048e0798ec9ae7dceb8",
    importpath = "gopkg.in/tomb.v1",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "51d6538a90f86fe93ac480b35f37b2be17fef232",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "io_k8s_client_go",
    commit = "d082d5923d3cc0bfbb066ee5fbdea3d0ca79acf8",
    importpath = "k8s.io/client-go",
)

go_repository(
    name = "io_k8s_code_generator",
    commit = "639c964206c28ac3859cf36f212c24775616884a",
    importpath = "k8s.io/code-generator",
)

go_repository(
    name = "io_k8s_gengo",
    commit = "e17681d19d3ac4837a019ece36c2a0ec31ffe985",
    importpath = "k8s.io/gengo",
)

go_repository(
    name = "io_k8s_klog",
    commit = "e531227889390a39d9533dde61f590fe9f4b0035",
    importpath = "k8s.io/klog",
)

go_repository(
    name = "io_k8s_kube_openapi",
    commit = "a01b7d5d6c2258c80a4a10070f3dee9cd575d9c7",
    importpath = "k8s.io/kube-openapi",
)

go_repository(
    name = "io_k8s_sigs_controller_runtime",
    commit = "f6f0bc9611363b43664d08fb097ab13243ef621d",
    importpath = "sigs.k8s.io/controller-runtime",
    build_extra_args = ["-exclude=vendor"],
)

go_repository(
    name = "io_k8s_sigs_controller_tools",
    commit = "950a0e88e4effb864253b3c7504b326cc83b9d11",
    importpath = "sigs.k8s.io/controller-tools",
)

go_repository(
    name = "io_k8s_sigs_testing_frameworks",
    commit = "d348cb12705b516376e0c323bacca72b00a78425",
    importpath = "sigs.k8s.io/testing_frameworks",
)

go_repository(
    name = "org_golang_google_appengine",
    commit = "54a98f90d1c46b7731eb8fb305d2a321c30ef610",
    importpath = "google.golang.org/appengine",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "cbcb750295291b33242907a04be40e80801d0cfc",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "a4d6f7feada510cc50e69a37b484cb0fdc6b7876",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "9f3314589c9a9136388751d9adae6b0ed400978a",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "a5b02f93d862f065920dd6a40dddc66b60d0dec4",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    commit = "342b2e1fbaa52c93f31447ad2c6abc048c63e475",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "org_golang_x_time",
    commit = "9d24e82272b4f38b78bc8cff74fa936d31ccd8ef",
    importpath = "golang.org/x/time",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "99f201b6807eb28f750a1966316bb0d4417b6020",
    importpath = "golang.org/x/tools",
)

go_repository(
    name = "org_uber_go_atomic",
    commit = "df976f2515e274675050de7b3f42545de80594fd",
    importpath = "go.uber.org/atomic",
)

go_repository(
    name = "org_uber_go_multierr",
    commit = "3c4937480c32f4c13a875a1829af76c98ca3d40a",
    importpath = "go.uber.org/multierr",
)

go_repository(
    name = "org_uber_go_zap",
    commit = "27376062155ad36be76b0f12cf1572a221d3a48c",
    importpath = "go.uber.org/zap",
)
