const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// 获取所有 Vue 文件
function getAllVueFiles(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  
  list.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    
    if (stat && stat.isDirectory()) {
      results = results.concat(getAllVueFiles(filePath));
    } else if (path.extname(filePath) === '.vue') {
      results.push(filePath);
    }
  });
  
  return results;
}

// 检查 Vue 文件的语法
function checkVueFile(filePath) {
  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    
    // 检查未闭合的标签
    const openTags = [];
    const lines = content.split('\n');
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      
      // 匹配开始标签，但排除自闭合标签和注释
      const openTagMatches = line.match(/<([a-zA-Z][a-zA-Z0-9-]*)(?:\s+[^>]*)?(?<!\/|!--)>/g);
      if (openTagMatches) {
        openTagMatches.forEach(match => {
          const tagName = match.match(/<([a-zA-Z][a-zA-Z0-9-]*)/)[1];
          openTags.push({ tag: tagName, line: i + 1 });
        });
      }
      
      // 匹配结束标签
      const closeTagMatches = line.match(/<\/([a-zA-Z][a-zA-Z0-9-]*)>/g);
      if (closeTagMatches) {
        closeTagMatches.forEach(match => {
          const tagName = match.match(/<\/([a-zA-Z][a-zA-Z0-9-]*)>/)[1];
          
          // 查找最近的匹配开始标签
          let found = false;
          for (let j = openTags.length - 1; j >= 0; j--) {
            if (openTags[j].tag === tagName) {
              openTags.splice(j, 1);
              found = true;
              break;
            }
          }
          
          if (!found) {
            console.error(`${filePath}:${i + 1} - 多余的结束标签 </${tagName}>`);
          }
        });
      }
      
      // 检查自闭合标签
      const selfClosingMatches = line.match(/<([a-zA-Z][a-zA-Z0-9-]*)(?:\s+[^>]*)?\s*\/>/g);
      if (selfClosingMatches) {
        // 自闭合标签不需要结束标签，所以不做处理
      }
    }
    
    // 输出未闭合的标签
    if (openTags.length > 0) {
      console.error(`${filePath} - 未闭合的标签:`);
      openTags.forEach(tag => {
        console.error(`  行 ${tag.line}: <${tag.tag}>`);
      });
      return false;
    }
    
    return true;
  } catch (error) {
    console.error(`${filePath} - 错误: ${error.message}`);
    return false;
  }
}

// 主函数
function main() {
  const srcDir = path.join(__dirname, 'src');
  const vueFiles = getAllVueFiles(srcDir);
  
  console.log(`找到 ${vueFiles.length} 个 Vue 文件`);
  
  let hasErrors = false;
  
  vueFiles.forEach(file => {
    const relativePath = path.relative(__dirname, file);
    process.stdout.write(`检查 ${relativePath} ... `);
    
    const isValid = checkVueFile(file);
    
    if (isValid) {
      console.log('✓ 有效');
    } else {
      console.log('❌ 无效');
      hasErrors = true;
    }
  });
  
  if (hasErrors) {
    process.exit(1);
  }
}

main();
