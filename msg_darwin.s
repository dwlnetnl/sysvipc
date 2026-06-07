#include "textflag.h"

TEXT libSystem_msgget_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_msgget(SB)
GLOBL ·libSystem_msgget_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_msgget_trampoline_addr(SB)/8, $libSystem_msgget_trampoline<>(SB)

TEXT libSystem_msgctl_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_msgctl(SB)
GLOBL ·libSystem_msgctl_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_msgctl_trampoline_addr(SB)/8, $libSystem_msgctl_trampoline<>(SB)

TEXT libSystem_msgsnd_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_msgsnd(SB)
GLOBL ·libSystem_msgsnd_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_msgsnd_trampoline_addr(SB)/8, $libSystem_msgsnd_trampoline<>(SB)

TEXT libSystem_msgrcv_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_msgrcv(SB)
GLOBL ·libSystem_msgrcv_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_msgrcv_trampoline_addr(SB)/8, $libSystem_msgrcv_trampoline<>(SB)
