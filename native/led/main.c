#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdint.h>
#include <famap.h>

#define GPIO_PHY_BASE 0xC000A000
#define GPIO_PHY_SIZE 0xC0

struct led_pinmap
{
    int module;
    int bit;
    int value;
};

struct led_pinmap pinmaps[] = {
    {1, 31, 1},
    {2, 0, 1},
    {2, 1, 1}
};

static mmap_alloc_t  io_map;
static volatile uint32_t *io_port;

void  io_init( void )
{
	io_port = (uint32_t *)fa_mmap_alloc( &io_map, GPIO_PHY_BASE , GPIO_PHY_SIZE );
}

void  io_exit( void )
{
	fa_mmap_free( &io_map );
}

/*
GPIO의 ALTERNATE FUNCTION을 결정한다.

module: 
0: GPIOA
1: GPIOB
2: GPIOC

bit:
GPIO 번호

value:
0 : GPIO 
1 : ALT Function1
2 : ALT Function2 
3 : Reserved

예)
alt(1, 1, 0) => GPIOB1를 GPIO로 사용한다.
*/
void alt(int module, int bit, int value)
{
    int alt_addr;
    alt_addr = 0x10 * module + 8;

    uint32_t low_mask = (1 << ((bit % 16) * 2));
    uint32_t high_mask = (1 << ((bit % 16) * 2 + 1));

    // 값을 입력하기 전에 우선 해당 bit를 clear한다.
    io_port[alt_addr + (bit/16)] &= ~(high_mask | low_mask);

    if (value & 0x01)
    {
        io_port[alt_addr + (bit/16)] |= low_mask;
    }

    if (value & 0x02)
    {
        io_port[alt_addr + (bit/16)] |= high_mask;
    }
}

/*
GPIO의 INPUT/OUTPUT을 결정한다.

module: 
0: GPIOA
1: GPIOB
2: GPIOC

bit:
GPIO 번호

value:
0: INPUT
1: OUTPUT

예)
dir(1, 1, 0) => GPIOB1를 INPUT으로 사용한다.
dir(2, 0, 0) => GPIOC0를 OUTPUT으로 사용한다.
*/

void dir(int module, int bit, int value)
{
    int dir_addr;
    dir_addr = 0x10 * module + 1;

    if (value)
    {
        io_port[dir_addr] |= (1 << bit);
    }
    else
    {
        io_port[dir_addr] &= ~(1 << bit);
    }
}

/*
GPIO에 값을 출력한다. 정상적으로 동작하기 위해서는 아래의 조건이 맞아야한다.
- alt 함수를 사용하여 GPIO 모드로 설정되어 있어야 한다.
- dir 함수를 사용하여 OUTPUT 모드로 설정되어 있어야 한다.

module: 
0: GPIOA
1: GPIOB
2: GPIOC

bit:
GPIO 번호

value:
0: LOW
1: HIGH

예)
output(1, 1, 0) => GPIOB1에 LOW로 출력한다.
output(2, 0, 1) => GPIOC0에 HIGH로 출력한다.
*/
void output(int module, int bit, int value)
{
    int val_addr;

    val_addr = 0x10 * module;

    if (value)
    {
        io_port[val_addr] |= (1 << bit);
    }
    else
    {
        io_port[val_addr] &= ~(1 << bit);
    }
}

void led(int index, int value)
{

}
int main(int argc, char ** argv)
{
    int i;
    int mode = 0;

    if (argc == 2)
    {
        mode = argv[1][0] - '0';
        if (mode < 0 || mode > 4)
        {
            mode = 0;
        }
    }

    //mode:
    // 0x00  D9 켜기
    // 0x01  D8 켜기
    // 0x02  D8 켜기, D10 블링킹

    io_init();

    for (i = 0; i < sizeof(pinmaps)/sizeof(struct led_pinmap); i++)
    {
        alt(pinmaps[i].module, pinmaps[i].bit, 0);
        dir(pinmaps[i].module, pinmaps[i].bit, 1);
        pinmaps[i].value = 1;
        output(pinmaps[i].module, pinmaps[i].bit, pinmaps[i].value);
    }
    
    // Init
    switch(mode)
    {
        case 0:
            pinmaps[2].value = 0;
            output(pinmaps[2].module, pinmaps[2].bit, pinmaps[2].value);
            break;
        case 1:
            pinmaps[0].value = 0;
            output(pinmaps[0].module, pinmaps[0].bit, pinmaps[0].value);
            break;
        case 2:
            pinmaps[0].value = 0;
            output(pinmaps[0].module, pinmaps[0].bit, pinmaps[0].value);
            break;
    }

    // Blinking
    while (1)
    {
        switch(mode)
        {
            case 1:
                break;
            case 2:
                pinmaps[1].value = !pinmaps[1].value;
                output(pinmaps[1].module, pinmaps[1].bit, pinmaps[1].value);
                break;
            case 3:
                break;
        }
        usleep(500 * 1000);
    }

    io_exit();

    return 0;
}
