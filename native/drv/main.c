/*
 * file     main.c
 * Date     2013-06-10
 * modify	2013-08-08
 * Author   khkraining@falinux.com
 *
 * Copyright(C) 2013, falinux : Embedded Solution Division
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <signal.h>

#include <sys/ioctl.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/poll.h>

#include <linux/limits.h>
#include <linux/kdev_t.h>

#include <termios.h>
#include <unistd.h>

#include <webconn.h>

int dev_fd;

/*
 * Device Open
 */
int dev_open(
			char *fname,
			unsigned char major,
			unsigned char minor
			)
{
    int dev;

    dev = open( fname, O_RDWR | O_NDELAY );
    if (dev < 0) {
        if(access(fname, F_OK) == 0) {
            unlink(fname);
        }
        mknod(fname, (S_IRWXU|S_IRWXG|S_IFCHR), MKDEV(major,minor));

        dev = open(fname, O_RDWR|O_NDELAY);
        if (dev < 0) {
            printf( "Device OPEN FAIL %s\n", fname );
            return -1;
        }
    }
    return dev;
}

/*
 *	sub list
 */
void sub_list( void )
{
    printf( "\n************************\n");
    printf( "# [1] IOCTL WRITE\n");
    printf( "# [2] IOCTL READ\n");
    printf( "# [q] Exit\n");
    printf( "--------------------------\n");
    printf("\rCommand ? " );
}

/*
 * ioctl test
 */
void ioctl_test(void)
{
    unsigned int key;
    unsigned int val;

	sub_list();

    while (1) {
        key = getchar();
        if (key == 'q')
			break;

        switch (key) {
        case '1' :  val = 1;
					ioctl(dev_fd, IOCTL_WRITE, &val);
					sub_list();
					break;
        case '2' :  ioctl( dev_fd, IOCTL_READ, &val);
					printf("app-->val = %d\n");
					sub_list();
					break;
        }
    }
}

/*
 * Command List
 */
void command_list(void)
{
    printf("\r------------------------------------------\n");
    printf("\r 1: ioctl_test\n");
    printf("\r q: Exit\n" );
    printf("\r------------------------------------------\n");
    printf("\r select : ");
}

int main (void)
{
	unsigned int mainkey;

	char buf[128] = {0,};

	printf("Program Start...\n");
	

	dev_fd = dev_open("/dev/webconn", WEBCONN_MAJOR, 0);

	write(dev_fd, "khkrai", 128);

	read(dev_fd, &buf, 128);

	printf("\n\nbuf=%s\n", buf);

	printf("%x\n", IOCTL_WRITE);
	printf("%x\n", IOCTL_READ);
	ioctl(dev_fd, IOCTL_WRITE, 1);
	ioctl(dev_fd, IOCTL_READ, 0);


	close(dev_fd);
		return 0;
}
