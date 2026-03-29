#include "textflag.h"

#define ZERO(r) \
    XORQ r, r

TEXT ·LowerBound(SB), NOSPLIT, $0
    MOVQ slice_ptr+0(FP), AX
    MOVQ slice_len+8(FP), DX
    MOVQ index+24(FP), BX

    ZERO(R8)
    MOVQ DX, R9

loop:
    CMPQ R8, R9
    JGE done

    MOVQ R9, R10
    SUBQ R8, R10
    SHRQ $1, R10
    ADDQ R8, R10

    MOVQ (AX)(R10*8), R11

    MOVQ R10, R12
    ADDQ $1, R12

    CMPQ R11, BX
    CMOVQLT R12, R8
    CMOVQGE R10, R9

    JMP loop

done:
    MOVQ R8, res+32(FP)
    RET
