from HTMLParser import HTMLParser
from smartypants import smartypants
import sqlite3
import re

def _to_smart(verse):
    verse = verse.replace(",`",", '")
    verse = verse.replace("`","'")
    out = smartypants(verse)
    parser = HTMLParser()
    out = parser.unescape(out)

    return out

con = sqlite3.connect("bible_old.db")
versions = ["NIV", "ESV", "HCSB"]

cur = con.cursor()
cur.execute("SELECT DISTINCT BOOK FROM BIBLE WHERE VERSION=?", ("NIV",))
books = cur.fetchall()
booklist = []
for b in books:
    booklist.append(b[0])
    
for version in versions:
    #query each verse in each book
    #convert to smartypants
    #HTMLParser.unescape
    for book in booklist:
        file = open(version + "_" + book + "_sql_script.sql", "w")
        file.write("BEGIN TRANSACTION;\n")
        cur.execute("SELECT * FROM BIBLE WHERE BOOK=? AND VERSION=?", (book, version,))
        verses = cur.fetchall()
        for v in verses:
            chap = v[2]
            verseno = v[3]
            text = v[4]
            text = _to_smart(text)
            file.write("INSERT INTO BIBLE VALUES\n")
            insert_str = ("('%s', '%s', %d, %d, '%s');\n") % (version, book, chap, verseno, text)
            file.write(insert_str.encode('utf8'))
        file.write("COMMIT;\n")
        file.close()