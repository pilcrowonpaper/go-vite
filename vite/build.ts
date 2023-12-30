import { build } from "vite";

import fs from "fs/promises";
import path from "path";

import type { RollupOutput } from "rollup";

async function getAllHTMLFilePaths(dirPath: string): Promise<string[]> {
  const contentNames = await fs.readdir(dirPath);
  const result: string[] = [];
  for (const contentName of contentNames) {
    const contentPath = path.join(dirPath, contentName);
    if (contentName === "node_modules" || contentName.startsWith(".")) {
      continue;
    }
    const stat = await fs.stat(contentPath);
    if (stat.isDirectory()) {
      const nested = await getAllHTMLFilePaths(contentPath);
      result.push(...nested);
    } else if (contentName.endsWith(".html")) {
      result.push(contentPath);
    }
  }
  return result;
}

const buildResult = (await build({
  logLevel: "error",
  build: {
    write: false,
    rollupOptions: {
      input: await getAllHTMLFilePaths(process.cwd()),
    },
  },
})) as RollupOutput;

await fs.rm("vite/.html", {
  recursive: true,
  force: true,
});

await fs.rm("vite/.assets", {
  recursive: true,
  force: true,
});

await fs.mkdir("vite/.html", {
  recursive: true,
});

await fs.mkdir("vite/.assets", {
  recursive: true,
});

for (const outputFile of buildResult.output) {
  if (outputFile.type === "chunk") {
    await fs.writeFile(
      path.join("vite/.assets", outputFile.fileName.split("/").at(-1)!),
      outputFile.code
    );
  }
  if (outputFile.type === "asset" && !outputFile.fileName.endsWith(".html")) {
    await fs.writeFile(
      path.join("vite/.assets", outputFile.fileName.split("/").at(-1)!),
      outputFile.source
    );
  }
  if (outputFile.type === "asset" && outputFile.fileName.endsWith(".html")) {
    const id = Buffer.from(outputFile.fileName).toString("base64url");
    await fs.writeFile(
      path.join("vite/.html", `${id}.html`),
      outputFile.source
    );
  }
}
