# ====================================================================================
# Setup Project
PROJECT_NAME := provider-dummy
PROJECT_REPO := github.com/upbound/$(PROJECT_NAME)

PLATFORMS ?= linux_amd64 linux_arm64
-include build/makelib/common.mk

# ====================================================================================
# Setup Output

-include build/makelib/output.mk

# ====================================================================================
# Setup Go

NPROCS ?= 1
GO_TEST_PARALLEL := $(shell echo $$(( $(NPROCS) / 2 )))
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/provider
GO_LDFLAGS += -X $(GO_PROJECT)/internal/version.Version=$(VERSION)
GO_SUBDIRS += cmd internal apis
GO111MODULE = on
-include build/makelib/golang.mk

# ====================================================================================
# Setup Kubernetes tools

UP_VERSION = v0.16.1
-include build/makelib/k8s_tools.mk

# ====================================================================================
# Setup Images

IMAGES = provider-dummy
-include build/makelib/imagelight.mk

# ====================================================================================
# Setup XPKG

XPKG_REG_ORGS ?= xpkg.upbound.io/upbound
# NOTE(hasheddan): skip promoting on xpkg.upbound.io as channel tags are
# inferred.
XPKG_REG_ORGS_NO_PROMOTE ?= xpkg.upbound.io/upbound
XPKGS = provider-dummy
-include build/makelib/xpkg.mk

# NOTE(hasheddan): we force image building to happen prior to xpkg build so that
# we ensure image is present in daemon.
xpkg.build.provider-dummy: do.build.images

fallthrough: submodules
	@echo Initial setup complete. Running make again . . .
	@make

# integration tests
e2e.run: test-integration

# Run integration tests.
test-integration: $(KIND) $(KUBECTL) $(UP) $(HELM3)
	@$(INFO) running integration tests using kind $(KIND_VERSION)
	@KIND_NODE_IMAGE_TAG=${KIND_NODE_IMAGE_TAG} $(ROOT_DIR)/cluster/local/integration_tests.sh || $(FAIL)
	@$(OK) integration tests passed

# Update the submodules, such as the common build scripts.
submodules:
	@git submodule sync
	@git submodule update --init --recursive

# NOTE(hasheddan): the build submodule currently overrides XDG_CACHE_HOME in
# order to force the Helm 3 to use the .work/helm directory. This causes Go on
# Linux machines to use that directory as the build cache as well. We should
# adjust this behavior in the build submodule because it is also causing Linux
# users to duplicate their build cache, but for now we just make it easier to
# identify its location in CI so that we cache between builds.
go.cachedir:
	@go env GOCACHE

# NOTE(hasheddan): we must ensure up is installed in tool cache prior to build
# as including the k8s_tools machinery prior to the xpkg machinery sets UP to
# point to tool cache.
build.init: $(UP)

# This is for running out-of-cluster locally, and is for convenience. Running
# this make target will print out the command which was used. For more control,
# try running the binary directly with different arguments.
run: go.build
	@$(INFO) Running Crossplane locally out-of-cluster . . .
	@# To see other arguments that can be provided, run the command with --help instead
	DO_NOTHING=false $(GO_OUT_DIR)/provider --debug

dev: $(KUBECTL)
	@$(INFO) Deploying dummy server
	@$(KIND) load docker-image $(BUILD_REGISTRY)/server-dummy-$(HOSTARCH)
	@cat cluster/server-deployment.yaml | sed 's|muvaf/server-dummy:latest|$(BUILD_REGISTRY)/server-dummy-$(HOSTARCH)|' | $(KUBECTL) apply -f -
	@$(OK) Deploying dummy server
	@$(INFO) Installing Provider Dummy CRDs
	@$(KUBECTL) apply -R -f package/crds
	@$(OK) Installing Provider Dummy CRDs
	@$(INFO) Creating ProviderConfig
	@$(KUBECTL) apply -f examples/providerconfig/local.yaml
	@$(OK) Creating ProviderConfig
	@$(INFO) Port-forwarding to dummy server
	@$(KUBECTL) port-forward svc/server-dummy 8080:80 &
	@$(OK) Port-forwarding to dummy server
	@$(INFO) Starting Provider Dummy controllers
	@DO_NOTHING=false $(GO) run cmd/provider/main.go --debug

dev-clean: $(KIND) $(KUBECTL)
	@$(INFO) Deleting dummy server
	@$(KUBECTL) delete -f cluster/server-deployment.yaml
	@$(OK) Deleting dummy server
	@$(INFO) Deleting Objects
	@$(KUBECTL) delete dummy --all
	@$(OK) Deleting Objects
	@$(INFO) Deleting CRDs
	@$(KUBECTL) delete -f package/crds
	@$(OK) Deleting CRDs

.PHONY: submodules fallthrough test-integration run dev dev-clean

define CROSSPLANE_MAKE_HELP
Crossplane Targets:
    submodules            Update the submodules, such as the common build scripts.
    run                   Run crossplane locally, out-of-cluster. Useful for development.

endef
# The reason CROSSPLANE_MAKE_HELP is used instead of CROSSPLANE_HELP is because the crossplane
# binary will try to use CROSSPLANE_HELP if it is set, and this is for something different.
export CROSSPLANE_MAKE_HELP

crossplane.help:
	@echo "$$CROSSPLANE_MAKE_HELP"

help-special: crossplane.help

.PHONY: crossplane.help help-special
