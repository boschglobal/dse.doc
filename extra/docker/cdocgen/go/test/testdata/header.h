/**
 *  MyHeader.h
 *  ==========
 *
 *  This is a header file for a custom library.
 */

/**
 *  Valid H1
 *  --------
 *
 *  Valid paragraph. Module level doc.
 */

/**
 *  # Not Valid H1
 *  
 *  Valid paragraph.
 */

/**
 *  # Not Valid H2
 *
 *  Valid paragraph.
 */

/*
 *  Not Valid H1 1 star
 *  --------
 *
 *  Valid paragraph.
 */

/***
 *  Not Valid H1 3 star
 *  --------
 *
 *  Valid paragraph.
 */

/**
 *  Not Valid H1
 * 
 *  Valid paragraph.
 */

#ifndef MY_HEADER_H
#define MY_HEADER_H

/**
 *  Module Level Doc
 *  ================
 *
 *  This library provides various utility functions and data structures.
 *  It can be used for mathematical calculations, string manipulation, and more.
 *  For functions that operate on specific data types, see their respective sections below.
 */




/**
 *  Valid paragraph
 *  
 *  Valid H1
 *  --------
 *
 *  Valid paragraph.
 */

/**
 *  MyStruct
 *
 *  This is a typedef called MyStruct.
 *  Example
 *  -------
 *       Extract from header file __only__!!!!
 */
typedef struct MyStruct {
    int x;
    int y;
} MyStruct;


int myFunction(int param1, int param2);

/**
 *  anotherFunction
 *
 *  ...
 * Example
 * -----
 *
 * Extract from all files.
 */
int anotherFunction();


/**
 *  invalidFunction
 *
 *  
 * No description provided.
 */
void invalidFunction();

int missingFunction();

#endif /* MY_HEADER_H */
