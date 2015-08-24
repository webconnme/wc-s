BASE_PATH=/nfs/webconn

for MODULE in src/wc-s/*
do
	MODULE_NAME=$(basename ${MODULE})
	mkdir -p ${BASE_PATH}/${MODULE_NAME}
	
	cp bin/linux_arm/${MODULE_NAME} ${BASE_PATH}/${MODULE_NAME}/
	if [ -d files/${MODULE_NAME} ]
	then
		cp -Rf files/${MODULE_NAME}/*  ${BASE_PATH}/${MODULE_NAME}/
	fi
done

