include ../../makefile.mk

js:
	$(GO) run $(GOGIO) -target js -o ./wasm .

android:
	PATH=$(ANDROID_JAVA_ROOT) java -version
	PATH=$(ANDROID_JAVA_ROOT):$(PATH) ANDROID_SDK_ROOT=$(ANDROID_SDK_ROOT) $(GO) run $(GOGIO) -target android -o ./android.apk .

ios:
	$(GO) run $(GOGIO) -target ios -o ./ios.app .

windows:
	$(GO) run $(GOGIO) -target windows -arch amd64 -o ./windows_amd64.exe .

macos:
	$(GO) run $(GOGIO) -target macos -arch arm64 -o ./macos_arm64.app .
	$(GO) run $(GOGIO) -target macos -arch amd64 -o ./macos_intel.app .

all: js android ios windows macos