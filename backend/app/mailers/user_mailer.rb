
class UserMailer < ApplicationMailer
  default from: 'yammanuru.tejaswini@vegrow.in' 

  def welcome_email(user)
    @user = user
    mail(to: @user.email, subject: "🎉 Welcome to the Library Portal!")
  end
end
