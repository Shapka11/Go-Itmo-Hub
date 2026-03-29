#include "textflag.h"

#define ZERO(r) \
    XORQ r, r

TEXT ·Sum(SB), NOSPLIT, $0
    MOVQ slice_ptr+0(FP), AX
    MOVQ slice_len+8(FP), DX
    ZERO(R8)

loop:
    CMPQ DX, $0
    JE done

    MOVLQSX (AX), R9
    ADDQ R9, R8
    ADDQ $4, AX
    SUBQ $1, DX
    JMP loop

done:
    MOVQ R8, res+24(FP)
    RET
