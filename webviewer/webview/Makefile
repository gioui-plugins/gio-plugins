include ../../makefile.mk

$(eval $(call generate_java,sys))
$(eval $(call generate_inkwasm,sys))

android: java_sys

js: inkwasm_sys

gen: android js