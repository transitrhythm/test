// test_sqlite3.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <sqlite3.h>

int main()
{
	sqlite3* db;
	std::cout << "Hello My Sqlite3!\n";
	int fd = sqlite3_open("mydb.db", &db);

	if (fd == SQLITE_OK) {
		std::cout << "Success opening the database.\n";
	}
	else {
		std::cerr << "Error:\n";
		std::cerr << sqlite3_errmsg(db) << '\n';
		exit(1);
	}
	sqlite3_close(db);
}

// Run program: Ctrl + F5 or Debug > Start Without Debugging menu
// Debug program: F5 or Debug > Start Debugging menu

// Tips for Getting Started: 
//   1. Use the Solution Explorer window to add/manage files
//   2. Use the Team Explorer window to connect to source control
//   3. Use the Output window to see build output and other messages
//   4. Use the Error List window to view errors
//   5. Go to Project > Add New Item to create new code files, or Project > Add Existing Item to add existing code files to the project
//   6. In the future, to open this project again, go to File > Open > Project and select the .sln file
