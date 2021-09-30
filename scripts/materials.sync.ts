import puppeteer from "puppeteer";
import {writeFileSync} from "fs";

(async () => {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();

    await page.goto("https://wiki.biligame.com/ys/%E6%9D%90%E6%96%99%E4%B8%80%E8%A7%88", {});

    const ret = await page.evaluate(() => {
        const ret: {
            [name: string]: {
                Icon?: string
                DropWeekdays?: number[],
            }
        } = {}


        document.querySelectorAll("table#CardSelectTr tbody > tr[data-param1]").forEach(($tr, i) => {
            if (i === 0) {
                return
            }

            const $firstTd = $tr.querySelectorAll("td").item(0)

            if (!$firstTd) {
                return
            }

            const name = $firstTd.querySelector("a")?.getAttribute("title") || ""
            const icon = $firstTd.querySelector("a > img")?.getAttribute("src") || ""

            const matched = $tr.getAttribute("data-param2")?.match(/时间：周(.)\/周(.)\/周(.)/)

            ret[name] = {
                Icon: icon,
                DropWeekdays: matched ? matched.slice(1).map((v) => "日一二三四五六".indexOf(v)).sort() : undefined,
            }
        })

        return ret
    });

    await browser.close();

    writeFileSync("./genshindb/materials.json", JSON.stringify(ret, null, 2))
})();
