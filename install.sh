BASE_PATH=/nfs/webconn/pkg

mkdir -p ${BASE_PATH}/app/webconn/bin
for MODULE in src/wc-s/app/*
do
	MODULE_NAME=$(basename ${MODULE})
	cp bin/linux_arm/${MODULE_NAME} ${BASE_PATH}/app/webconn/bin
done

mkdir -p ${BASE_PATH}/env/falinux/bin
for MODULE in src/wc-s/env/*
do
	MODULE_NAME=$(basename ${MODULE})
	cp bin/linux_arm/${MODULE_NAME} ${BASE_PATH}/env/falinux/bin
done

cp -R files/*  ${BASE_PATH}/
