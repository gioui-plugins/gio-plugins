include ../makefile.mk

$(eval $(call generate_java,auth,"./vendors/android/target/*"))
$(eval $(call generate_pom,auth,"./vendors/android"))

android: pom_auth java_auth

gen: android