from selenium import webdriver
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.service import Service
import time
import random
import pytz
import sys
from datetime import datetime

run_time = 1800
log_name = 'spider.log'
driver_path = Service("/chrome_soft/chromedriver_v100")
chrome_options = Options()
prefs = { 
    'profile.default_content_setting_values':
    {
        'notifications': 2
    }
}
chrome_options.add_experimental_option('prefs', prefs)
chrome_options.add_argument('--headless')
chrome_options.add_argument('--first-run')
chrome_options.add_argument('--start-maximized')
chrome_options.add_argument("--disable-component-update")
chrome_options.add_argument("--disable-application-cache")
chrome_options.add_argument('--disable-gpu')
chrome_options.add_argument('--ignore-ssl-error')
chrome_options.add_argument('--ignore-certificate-errors')
chrome_options.add_argument('--no-sandbox')
chrome_options.add_argument('--proxy-server=socks://172.19.0.3:11156')

words = []
with open('words.txt', 'r') as f:
    words = f.read().split('\n')
words_len = len(words)

def amazon(driver):
    driver.get('http://www.amazon.com')
    i = 0
    WebDriverWait(driver, 10).until(
            EC.visibility_of_element_located((By.NAME, "field-keywords")))
    while i <= 5:
        KW_num = random.randint(0, words_len)
        time.sleep(1)
        driver.find_element(by=By.NAME, value="field-keywords").send_keys(words[KW_num])
        time.sleep(1)
        driver.find_element(by=By.ID, value="nav-search-submit-button").click()
        time.sleep(1)
        if len(driver.find_elements(by=By.CLASS_NAME,value="s-image")) == 0:
            time.sleep(5)
        if len(driver.find_elements(by=By.CLASS_NAME,value="s-image")) == 0:
            time.sleep(10)
        if len(driver.find_elements(by=By.CLASS_NAME,value="s-image")) != 0:
            num2 = random.randint(0, len(driver.find_elements(by=By.CLASS_NAME,value="s-image")))
            driver.find_elements(by=By.CLASS_NAME,value="s-image")[num2].click()
            time.sleep(5)
            WebDriverWait(driver, 10).until(
                EC.visibility_of_element_located((By.NAME, "field-keywords")))
            driver.find_element(by=By.NAME, value="field-keywords").clear()
            time.sleep(5)
        i = i+1
    time.sleep(5)

def reddit(driver):
    driver.get('https://www.reddit.com')
    i = 0
    WebDriverWait(driver, 10).until(
        EC.visibility_of_element_located((By.CSS_SELECTOR, "#header-search-bar")))
    while i <= 5:
        KW_num = random.randint(0, words_len)
        time.sleep(1)
        driver.find_element(by=By.CSS_SELECTOR, value="#header-search-bar").send_keys(words[KW_num])
        time.sleep(1)
        driver.find_element(by=By.CSS_SELECTOR, value="#header-search-bar").submit()
        time.sleep(1)
        if len(driver.find_elements(by=By.CSS_SELECTOR,value=".y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE")) == 0:
            time.sleep(5)
        if len(driver.find_elements(by=By.CSS_SELECTOR,value=".y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE")) == 0:
            time.sleep(10)
        if len(driver.find_elements(by=By.CSS_SELECTOR,value=".y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE")) != 0:
            num2 = random.randint(0, len(driver.find_elements(by=By.CSS_SELECTOR,value=".y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE")))
            driver.find_elements(by=By.CSS_SELECTOR,value=".y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE")[num2].click()
            time.sleep(5)
            WebDriverWait(driver, 10).until(
               EC.visibility_of_element_located((By.CSS_SELECTOR, "#header-search-bar")))
            driver.find_element(by=By.CSS_SELECTOR, value="#header-search-bar").clear()
            time.sleep(5)
        i = i+1
    time.sleep(5)

def wiki(driver):
    driver.get('https://www.wikipedia.org')
    WebDriverWait(driver, 10).until(
        EC.visibility_of_element_located((By.XPATH, '//*[@id="www-wikipedia-org"]/div[2]/div[1]'))).click()
    i = 0
    while i<= 20:
        WebDriverWait(driver, 10).until(EC.visibility_of_element_located((By.XPATH, '//*[@id="n-randompage"]/a'))).click()
        time.sleep(5)
        i = i+1
    time.sleep(5)

def yahoo(driver):
    driver.get('https://www.yahoo.com')
    WebDriverWait(driver, 10).until(
        EC.visibility_of_element_located((By.ID, "ybar-sbq")))
    KW_num = random.randint(0, words_len)
    time.sleep(1)
    driver.find_element(by=By.ID, value="ybar-sbq").send_keys(words[KW_num])
    time.sleep(1)
    driver.find_element(by=By.ID, value="ybar-search").click()
    time.sleep(1)
    i = 0
    while i <= 10:
        WebDriverWait(driver, 10).until(
            EC.visibility_of_element_located((By.ID, "yschsp")))
        KW_num = random.randint(0, words_len)
        time.sleep(3)
        driver.find_element(by=By.ID, value="yschsp").clear()
        time.sleep(1)
        driver.find_element(by=By.ID, value="yschsp").send_keys(words[KW_num])
        time.sleep(1)
        driver.find_element(by=By.CLASS_NAME, value="sbb").click()
        time.sleep(1)
        i = i+1
    time.sleep(5)

def youtube(driver):
    driver.get('http://www.youtube.com')
    WebDriverWait(driver, 10).until(
        EC.visibility_of_element_located((By.NAME, "search_query"))).click()
    KW_num = random.randint(0, words_len)
    time.sleep(1)
    driver.find_element(by=By.NAME, value="search_query").send_keys(words[KW_num])
    time.sleep(1)
    WebDriverWait(driver, 10).until(EC.visibility_of_element_located(
        (By.ID, "search-icon-legacy"))).click()
    time.sleep(1)
    if len(driver.find_elements(by=By.ID,value="title-wrapper")) == 0:
        time.sleep(15)
    if len(driver.find_elements(by=By.ID,value="title-wrapper")) == 0:
        time.sleep(20)
    if len(driver.find_elements(by=By.ID,value="title-wrapper")) == 0:
        time.sleep(25)
    if len(driver.find_elements(by=By.ID,value="title-wrapper")) != 0:
        num2 = random.randint(0, len(driver.find_elements(by=By.ID,value="title-wrapper")))
        driver.find_elements(by=By.ID,value="title-wrapper")[num2].click()
        time.sleep(60)
    else:
        with open('spider.log','a') as f:
            now_time = datetime.now(pytz.timezone('Asia/Shanghai'))
            f.write(now_time.strftime("%Y-%m-%d_%H:%M:%S")+":[youtube]timeout \n")

if __name__ == "__main__":
    website = sys.argv[1]
    start_time = datetime.now(pytz.timezone('Asia/Shanghai'))
    with open('spider.log','a') as f:
        start_time = datetime.now(pytz.timezone('Asia/Shanghai'))
        f.write(start_time.strftime("%Y-%m-%d_%H:%M:%S")+":["+website+"]spdier start\n") 
    while True:
        time.sleep(5)
        driver = webdriver.Chrome(service=driver_path, options=chrome_options)
        try:
            eval(website)(driver)
        except Exception as e:
            error_string = str(e)
            if 'Stacktrace' in error_string:
                error_string = error_string.split('\n')[0]
            with open('spider.log','a') as f:
                now_time = datetime.now(pytz.timezone('Asia/Shanghai'))
                f.write(now_time.strftime("%Y-%m-%d_%H:%M:%S")+":["+website+"][ERROR] " + error_string +"\n") 
        driver.quit()
        now_time = datetime.now(pytz.timezone('Asia/Shanghai'))
        if (now_time - start_time).seconds >= run_time:
            break
    with open('spider.log','a') as f:
        now_time = datetime.now(pytz.timezone('Asia/Shanghai'))
        f.write(now_time.strftime("%Y-%m-%d_%H:%M:%S")+":["+website+"]spdier end\n") 
