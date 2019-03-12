#include<stdio.h>
#include<stdlib.h>
#include<time.h>

int getRandomTime(int max, int min){
    srand((unsigned)time(NULL));
    int a = min - 10;
    while(a < min){
        a = rand() % max;
    }
    
    return a;
}
