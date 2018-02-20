      store 1, #1
      sub A, B, A
      sub D, D, D
      jump LBL1
      add A, 15, A
LBL1: add A, B, B
      store B, #(D + 16)
      load #13, B
      store C, #14
