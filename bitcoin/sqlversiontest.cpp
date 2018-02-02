// With sqlite-dev package installed, compile this with:
// g++ -o sqlversiontest sqlversiontest.cpp -lsqlite3
//
#include <sqlite3.h>
#include <iostream>

using namespace std;

int main(int argc, char* argv[])
{
	cout << "SQLITE_VERSION_NUMBER=" << SQLITE_VERSION_NUMBER << "\n";
	cout << "sqlite3_libversion_number()=" << sqlite3_libversion_number() << "\n";

}
