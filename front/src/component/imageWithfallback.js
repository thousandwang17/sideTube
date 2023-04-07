/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-22 15:09:12
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-23 13:32:42
 * @FilePath: /sidetube/src/component/imageWithfallback.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Box } from "@mui/system";
import React, { useState, useEffect } from "react";
export default function ImageWithFallback({ primarySrc, fallbackSrc, style }) {
  const [src, setSrc] = useState(primarySrc);

  const handleError = () => {
    setSrc(fallbackSrc);
  };

  useEffect(() => {
    setSrc(primarySrc);
  }, [primarySrc]);

  return (
    <Box component="img" src={src} onError={handleError} style={{ ...style }} />
  );
}
