import json
import sys

from selenium import webdriver
from selenium.webdriver.chrome import service, options
from selenium.webdriver.common.by import By
from selenium.webdriver.support import wait


_ROADMAP_URL = "https://neetcode.io/roadmap"
_GROUPS_ORDER = [
    "Arrays & Hashing",
    "Stack",
    "Two Pointers",
    "Sliding Window",
    "Binary Search",
    "Linked List",
    "Trees",
    "Tries",
    "Backtracking",
    "Graphs",
    "Heap / Priority Queue",
    "Advanced Graphs",
    "Intervals",
    "Greedy",
    "1-D DP",
    "2-D DP",
    "Bit Manipulation",
    "Math & Geometry",
]


def main():
    opts = options.Options()
    opts.add_argument("--headless=new")
    driver = webdriver.Chrome(
        options=opts,
        service=service.Service(
            executable_path="./chromedriver",
            port=8787,
        ),
    )
    driver.get(_ROADMAP_URL)
    driver.implicitly_wait(0.5)

    result = []
    groups = driver.find_elements(By.CSS_SELECTOR, ".node-group label")
    print("loaded groups", [group.text for group in groups], file=sys.stderr)
    for group in groups:
        print(f"processing: {group.text}", file=sys.stderr)
        group.click()
        questions = []
        group_result = {
            "group_name": group.text,
            "questions": questions,
        }

        for question in driver.find_elements(By.CSS_SELECTOR, "table tbody tr"):
            cols = question.find_elements(By.CSS_SELECTOR, "td")
            problem_col, difficulty_col = cols[2], cols[3]
            wait.WebDriverWait(driver, timeout=2).until(
                lambda _: (problem_col.is_displayed() and difficulty_col.is_displayed())
            )

            name = problem_col.text
            difficulty = difficulty_col.find_element(By.CSS_SELECTOR, "b").text.lower()
            lc_link, free_link = "", ""
            for anchor in problem_col.find_elements(By.CSS_SELECTOR, "a"):
                link = anchor.get_attribute("href")
                if "leetcode" in link:
                    lc_link = link
                elif "neetcode" in link:
                    free_link = link

            assert name, group.text
            assert difficulty, group.text
            assert lc_link, group.text
            questions.append(
                {
                    "name": name,
                    "difficulty": difficulty,
                    "lc_link": lc_link,
                    "free_link": free_link,
                }
            )

        result.append(group_result)
        esc_btn = next(
            bt
            for bt in driver.find_elements(
                By.CSS_SELECTOR, ".my-sidebar .close-container button"
            )
            if bt.text == "ESC"
        )
        esc_btn.click()

    order_key = {group_name: i for i, group_name in enumerate(_GROUPS_ORDER)}
    result.sort(key=lambda e: order_key[e["group_name"]])

    print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
