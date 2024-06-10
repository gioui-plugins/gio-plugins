GO = "go"
GOGIO = "gioui.org/cmd/gogio"
ANDROID_SDK_ROOT = ""
ANDROID_JAVA_ROOT = ""
TEMP = ""

ifeq ($(shell uname -s), Darwin)
    ANDROID_SDK_ROOT = "$(HOME)/Library/Android/sdk/"
    ANDROID_JAVA_ROOT = /Library/Java/JavaVirtualMachines/zulu-11.jdk/Contents/Home/bin
    ANDROID_PLATFORM = $(ANDROID_SDK_ROOT)/platforms/$(shell ls $(ANDROID_SDK_ROOT)/platforms | sort -n | tail -n 1)
    TEMP = /tmp
endif

define generate_java
java_$(1):
	mkdir -p $(TEMP)/$(1)_android/classes
	PATH=$(ANDROID_JAVA_ROOT):$(PATH) javac -source 8 -target 8 -bootclasspath $(ANDROID_PLATFORM)/android.jar -cp $(2) -d $(TEMP)/$(1)_android/classes $(1)_android.java
	jar cf $(1)_android.jar -C $(TEMP)/$(1)_android/classes .
	rm -rf $(TEMP)/$(1)_android
endef

define generate_inkwasm
inkwasm_$(1):
	$(GO) run github.com/inkeliz/go_inkwasm@master generate .
	mv inkwasm_js.go $(1)_js_wasm.go
	mv inkwasm_js.s $(1)_js_wasm.s
	mv inkwasm_js.js $(1)_js.js
endef

CURDIR = $(shell pwd)

define generate_pom
pom_$(1):
	go run $(CURDIR)/../cmd/pomgen/pom.go $(2)
endef