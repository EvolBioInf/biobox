./sops -a 2  test.fasta > r1.txt
./sops -i -2 test.fasta > r2.txt
./sops -g -1 test.fasta > r3.txt
./sops -m sm.txt test.fasta > r4.txt
