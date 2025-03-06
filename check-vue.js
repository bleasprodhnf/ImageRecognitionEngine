const fs = require('fs');
const path = require('path');
const { parse } = require('@vue/compiler-sfc');

// 要检查的文件路径
const filePath = process.argv[2];

if (!filePath) {
  console.error('请提供要检查的 Vue 文件路径');
  process.exit(1);
}

try {
  const content = fs.readFileSync(filePath, 'utf-8');
  const { errors } = parse(content);
  
  if (errors && errors.length > 0) {
    console.error('❌ 语法错误:');
    errors.forEach(err => console.error(`- ${err}`));
    process.exit(1);
  } else {
    console.log(`✓ 文件 ${path.basename(filePath)} 语法有效`);
  }
} catch (error) {
  console.error(`❌ 错误: ${error.message}`);
  if (error.loc) {
    console.error(`  位置: 行 ${error.loc.start.line}, 列 ${error.loc.start.column}`);
  }
  process.exit(1);
}
