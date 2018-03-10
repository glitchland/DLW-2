      sub A, B, C
      sub C, B, A
      add A, B, C
      add C, B, A
      sub D, D, D
      add A, 15, C
      sub A, 15, C   
      store C, #(D + 16)
      store C, #1
      store B, #D
      load #13, B
      load #(C + 2), C
      load #D, D
      jumpz LBL1
      jumpz #C
      jumpz #(D + 16)
      jumpz 1
      jumpz -1      
      jump LBL1
      jump #C
      jump #(D + 16)
      jump 1
      jump -1  
LBL1: add A, B, B      
