export const mergeArraysSafe = (...arrays) => {
  // Filter arrays
  const filteredArrays = arrays.map(arr => Array.isArray(arr) ? arr : []);

  // Calculate total size
  const totalSize = filteredArrays.reduce((acc, arr) => acc + arr.length, 0);
  
  // Pre-allocate new array woth the size
  const mergedArray = new Array(totalSize);
  
  let offset = 0;
  for (const arr of filteredArrays) {
    // Copy elements
    for (let i = 0; i < arr.length; i++) {
      mergedArray[offset + i] = arr[i];
    }
    offset += arr.length;
  }
  
  return mergedArray;
};


const graphTransform = tables => tables.reduce((acc, table, index) => {
  // Add tables data
  acc.push({
    data: {
      id: index,
      name: table.tableName,
      columns: table.columns,
    },
  });

  // Add relations if exist
  const relationships = table.relationships || [];
  relationships.forEach(rel => {
    acc.push({
      data: {
        id: rel.conname,
        source: table.cableName,
        target: rel.relatedTableName,
      },
    });
  });

  return acc;
}, []);