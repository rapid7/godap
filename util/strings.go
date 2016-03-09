package util

func StringInSlice(str string, list []string) bool {
   for _, item := range list {
      if (str == item) {
         return true
      }
   }
   return false
}
