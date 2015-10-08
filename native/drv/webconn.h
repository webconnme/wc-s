/*
 * file     main.c
 * Date     2013-06-10
 * modify	2013-08-08
 * Author   khkraining@falinux.com
 *
 * Copyright(C) 2013, falinux : Embedded Solution Division
 */

#ifndef __DRV_BASIC_IOCTL_H_
#define __DRV_BASIC_IOCTL_H_

/*
 * IDTECK Major & Version Number
 */
#define WEBCONN_MAJOR	250
#define WEBCONN_VERSION	"ver 0.1"

/*
 * ioctl command
 */
#define	IOCTL_MAGIC	'I'

#define	IOCTL_WRITE	_IOW(IOCTL_MAGIC, 0, int )
#define	IOCTL_READ	_IOR(IOCTL_MAGIC, 1, int )

#define	IOCTL_MAXNR	2

#endif
