#include "textflag.h"

TEXT ·fileRead(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·fileWrite(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·fileSlice(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·writableWrite(SB), NOSPLIT, $0
  CallImport
  RET
