#include "textflag.h"

TEXT libSystem_semget_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_semget(SB)
GLOBL ·libSystem_semget_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_semget_trampoline_addr(SB)/8, $libSystem_semget_trampoline<>(SB)

TEXT libSystem_semop_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_semop(SB)
GLOBL ·libSystem_semop_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_semop_trampoline_addr(SB)/8, $libSystem_semop_trampoline<>(SB)

TEXT libSystem_semctl_trampoline<>(SB), NOSPLIT, $0-0
    JMP libSystem_semctl(SB)
GLOBL ·libSystem_semctl_trampoline_addr(SB), RODATA, $8
DATA  ·libSystem_semctl_trampoline_addr(SB)/8, $libSystem_semctl_trampoline<>(SB)
