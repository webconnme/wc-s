/**    
    @file     famap.c
    @date     2010/12/1
    @author   오재경 freefrug@falinux.com  FALinux.Co.,Ltd.
    @brief    * mmap 유틸리티이다.
              * tmmap.c 가 tlist.c 를 반드시 포함해야하는 문제점이 있어 이를 보완하기
                위해 작성되었다. (외부공개가 필요할때 사용하면 좋다.)
              * 사용법은 파일의 끝에 작성되어 있다.
              
              Ver 0.1.0
              
    @modify   
    @todo     
    @bug     
    @remark   
    @warning   tmmap.c  tmmap.h 와는 같이 사용하지 않는다.
*/
//----------------------------------------------------------------------------
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/poll.h>
#include <sys/mman.h> 
#include <famap.h>

//------------------------------------------------------------------------------
/** @brief   mmap 생성함수
    @param   ma           mmap_alloc_t 구조체 포인터
    @param   phys_base    물리주소
    @param   size         크기
    @return  사용할수 있는 가상주소
*///----------------------------------------------------------------------------
void  *fa_mmap_alloc( mmap_alloc_t *ma, unsigned long phys_base, unsigned long size )
{       
	int        dev_mem;
	int        base_ofs;
	void      *mmap_mem;
	long       base_addr;
	
	
	// 4K 정렬 주소로 변경한다.
	base_ofs   = phys_base & (PAGE_SIZE-1);
	base_addr = (long)phys_base & ~(PAGE_SIZE-1);
	phys_base = (unsigned long)base_addr;
	
	// 4K 단위의 메모리를 할당받는다.
	size = PAGE_SIZE * ( (size + base_ofs + (PAGE_SIZE-1))/(PAGE_SIZE) );
	
	dev_mem = open( "/dev/mem", O_RDWR|O_SYNC );
	if (0 > dev_mem)
	{
		printf( "open error /dev/mem\n" );
		return NULL;
	}

	// mmap 로 맵핑한다.
	mmap_mem = mmap( 0,                                // 커널에서 알아서 할당요청
                     size,                             // 할당 크기
                     PROT_READ|PROT_WRITE, MAP_SHARED, // 할당 속성
                     dev_mem,                          // 파일 핸들
                     phys_base );                      // 매핑 대상의 물리주소	
	
	
	if ( !mmap_mem )
	{
        close(dev_mem);
		printf( "mmap error !!!\n" );
		return NULL;
	}

	// 관리구조체를 채운다.
	ma->dev      = dev_mem;
	ma->phys     = phys_base;
	ma->size     = size;
	ma->virt     = mmap_mem;
	ma->base_ofs = base_ofs;

	return mmap_mem + base_ofs;
}

//------------------------------------------------------------------------------
/** @brief   mmap  포인터를 해제한다.
    @param   ma    mmap_alloc_t 구조체 포인터
*///----------------------------------------------------------------------------
void  fa_mmap_free( mmap_alloc_t *ma )
{
	munmap( ma->virt, ma->size );	
	close ( ma->dev );
}


/// 샘플 ------------------------------------------------------------------------
#if 0

#include <famap.h>

static mmap_alloc_t  io_map;
static volatile unsigned char *io_port;

void  io_init( void )
{
	io_port = (unsigned char *)fa_mmap_alloc( &io_map, PHY_BASE , PHY_SIZE );
	// io_port[0] = 0x00;
}

void  io_exit( void )
{
	fa_mmap_free( &io_map );
}

#endif
