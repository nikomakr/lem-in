#!/bin/bash

check() {
    local desc=$1
    local file=$2
    local expected=$3
    local result=$(go run . "$file" 2>/dev/null)
    if echo "$result" | grep -q "$expected"; then
        echo "✅ $result"
    else
        echo "❌ $desc → got: $result"
    fi
}

# 1. invalid number of ants
echo "0
##start
a 0 0
##end
b 1 1
a-b" > /tmp/t1.txt && check "invalid number of ants" /tmp/t1.txt "invalid number of ants"

# 2. no start room found
echo "2
a 0 0
##end
b 1 1
a-b" > /tmp/t2.txt && check "no start room found" /tmp/t2.txt "no start room found"

# 3. no end room found
echo "2
##start
a 0 0
b 1 1
a-b" > /tmp/t3.txt && check "no end room found" /tmp/t3.txt "no end room found"

# 4. multiple start rooms found
echo "2
##start
a 0 0
##start
b 1 1
##end
c 2 2
a-b
b-c" > /tmp/t4.txt && check "multiple start rooms found" /tmp/t4.txt "multiple start rooms found"

# 5. multiple end rooms found
echo "2
##start
a 0 0
##end
b 1 1
##end
c 2 2
a-b
b-c" > /tmp/t5.txt && check "multiple end rooms found" /tmp/t5.txt "multiple end rooms found"

# 6. duplicate room
echo "2
##start
a 0 0
a 1 1
##end
b 2 2
a-b" > /tmp/t6.txt && check "duplicate room" /tmp/t6.txt "duplicate room"

# 7. unknown room in link
echo "2
##start
a 0 0
##end
b 1 1
a-ghost" > /tmp/t7.txt && check "unknown room in link" /tmp/t7.txt "unknown room in link"

# 8. invalid room name
echo "1
##start
L1 0 0
##end
b 1 1
L1-b" > /tmp/t8.txt && check "invalid room name" /tmp/t8.txt "invalid room name"

# 9. invalid coordinates for room
echo "1
##start
a 1.5 2.5
##end
b 1 1
a-b" > /tmp/t9.txt && check "invalid coordinates for room" /tmp/t9.txt "invalid coordinates for room"

# 10. invalid tunnel format
echo "2
##start
a 0 0
##end
b 1 1
a-b-c" > /tmp/t10.txt && check "invalid tunnel format" /tmp/t10.txt "invalid tunnel format"

# 11. ##start declared but no room followed
echo "2
##start
a-b" > /tmp/t11.txt && check "##start declared but no room followed" /tmp/t11.txt "##start declared but no room followed"

# 12. ##end declared but no room followed
echo "2
##start
a 0 0
##end" > /tmp/t12.txt && check "##end declared but no room followed" /tmp/t12.txt "##end declared but no room followed"

# 13. no tunnels found
echo "2
##start
a 0 0
##end
b 1 1" > /tmp/t13.txt && check "no tunnels found" /tmp/t13.txt "no tunnels found"

# 14. file is empty
> /tmp/t14.txt && check "file is empty" /tmp/t14.txt "file is empty"

# 15. path is a directory not a file
check "path is a directory not a file" /tmp "path is a directory not a file"

# 16. could not open file
check "could not open file" /tmp/nonexistent.txt "could not open file"

# 17. file read error — simulated via a file with only whitespace
printf "   \n\n   \n" > /tmp/t17.txt && check "file contains no valid data" /tmp/t17.txt "file contains no valid data"

# 18. no path found between start and end
echo "2
##start
a 0 0
##end
b 1 1
c 2 2
a-c" > /tmp/t18.txt && check "no path found between start and end" /tmp/t18.txt "no path found between start and end"

echo ""
echo "Done"