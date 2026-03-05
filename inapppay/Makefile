include ../makefile.mk

$(eval $(call generate_java,inapppay,"./vendors/android/target/*"))
$(eval $(call generate_pom,inapppay,"./vendors/android"))

android: pom_inapppay java_inapppay

gen: android
