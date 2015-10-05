#
# 어플용 Makefile
# 

## 사용자 소스
TARGET          = app-gpio
SRCS            = main.c famap.c


OBJS			= $(SRCS:.c=.o)
CFLAGS			= -I. -Wall

CROSS_PREFIX    = arm-generic-linux-gnueabi-

CFLAGS         += -g

#---------------------------------------------------------------------
CC				=	$(CROSS_PREFIX)gcc
CXX			 	=	$(CROSS_PREFIX)g++
AR				=	$(CROSS_PREFIX)ar rc
AR2				=	$(CROSS_PREFIX)ranlib
RANLIB			=	$(CROSS_PREFIX)ranlib
LD				=	$(CROSS_PREFIX)ld
NM				=	$(CROSS_PREFIX)nm
STRIP			=	$(CROSS_PREFIX)strip
OBJCOPY		 	=	$(CROSS_PREFIX)objcopy
CP				=	cp
MV				=	mv
#--------------------------------------------------------------------

all : $(TARGET)
	chmod 777 $^

$(TARGET) : $(OBJS)
	@echo "Linking $@ ..."
	$(CC) $(LDFLAGS) $(OBJS) -o $@ $(LIBS) 

%.o:%.c
	@echo "Compiling $< ..."
	$(CC) -c $(CFLAGS) -o $@ $<

%.o:%.cc
	@echo "C++ compiling $< ..."
	$(CXX) -c $(CFLAGS) $(CXXFLAGS) -o $@ $<

dep :
	$(CC)	-M	$(INCLUDEDIRS)	$(SRCS) > .depend

clean:
	rm -f *.bak
	rm -f *.map
	rm -f *.o
	rm -f $(OBJS)
	rm -f $(TARGET) core

distclean: clean
	rm -rf .depend	

ifeq (.depend,$(wildcard .depend))
include .depend
endif

