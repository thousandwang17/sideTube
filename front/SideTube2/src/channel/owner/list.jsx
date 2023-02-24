/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-14 17:50:53
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-14 18:28:12
 * @FilePath: /sidetube/src/channel/videoUpload/upload.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from 'react';
import { Button } from '@mui/material';

export default function VideoOwnerList() {
 

    return (
        <div>
            <input
            accept="image/*"
            className=""
            style={{ display: 'none' }}
            id="raised-button-file"
            multiple
            type="file"
            />
            <label htmlFor="raised-button-file">
            <Button variant="raised" component="span" className="">
                Upload
            </Button>
            </label> 
        </div>
    );

}