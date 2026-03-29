#include "textflag.h"

TEXT ·Fibonacci(SB), NOSPLIT, $0
    MOVQ number+0(FP), CX
    MOVQ $0, AX
    MOVQ $1, DX

    CMPQ CX, $0
    JE zeroTakt

loop:
    MOVQ DX, BX
    ADDQ AX, DX
    MOVQ BX, AX

    SUBQ $1, CX

    CMPQ CX, $0
    JNE loop

    MOVQ BX, res+8(FP)
    RET

zeroTakt:
    MOVQ $0, res+8(FP)
    RET
