      sub A, B, A
      sub D, D, D
      jumpz LBL1
      add A, 15, A
      add #A, 15, A
      add B, #(D + 6), C
      add LBL1, C, #D
      store A, #(D + 16)
LBL1: add A, B, B
      store B, #(D + 16)
      load #13, B
      load #(A + 2), C
      store C, #14
      jumpz LBL1
      jumpz #C
      jumpz #(D + 16)
      jump LBL1
