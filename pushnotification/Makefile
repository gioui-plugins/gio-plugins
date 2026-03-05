include ../makefile.mk

$(eval $(call generate_java,pushnotification,"./vendors/android/target/*"))
$(eval $(call generate_pom,pushnotification,"./vendors/android"))

android: pom_pushnotification java_pushnotification

gen: android
